package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	followv1 "kitbook/api/proto/gen/follow/v1"
	"kitbook/feed/domain"
	"kitbook/feed/repository"
	"kitbook/pkg/logger"
	"sort"
	"sync"
	"time"
)

const (
	ArticleEventName = "article_event"
	// 粉丝数量阈值
	follweesThreshold = 100
)

type ArticleHandler struct {
	repo         repository.FeedEventRepository
	followClient followv1.FollowServiceClient
	l            logger.Logger
}

func NewArticleHandler(repo repository.FeedEventRepository, followClient followv1.FollowServiceClient, l logger.Logger) Handler {
	return &ArticleHandler{
		repo:         repo,
		followClient: followClient,
		l:            l,
	}
}

func (a *ArticleHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	followee, _ := ext.Get("followee").AsInt64()

	// 1. 判定粉丝数量
	resp, err := a.followClient.GetFollowStatics(ctx, &followv1.GetFollowStaticsRequest{
		Uid: followee,
	})
	if err != nil {
		return err
	}
	if resp.GetFollowers() > follweesThreshold {
		// 2. 拉模型
		return a.repo.CreatePullEvent(ctx, domain.FeedEvent{
			Type:  ArticleEventName,
			Uid:   followee,
			Ext:   ext,
			Ctime: time.Now(),
		})
	} else {
		// 2. 推模型
		// 2.1 先查出粉丝
		resp2, err := a.followClient.GetFollower(ctx, &followv1.GetFollowerRequest{
			Followee: followee,
		})
		if err != nil {
			return err
		}
		events := make([]domain.FeedEvent, 0, len(resp2.FollowRelations))
		for _, f := range resp2.FollowRelations {
			events = append(events, domain.FeedEvent{
				Type:  ArticleEventName,
				Uid:   f.Follower,
				Ext:   ext,
				Ctime: time.Now(),
			})
		}

		return a.repo.CreatePushEvents(ctx, events)
	}
}

func (a *ArticleHandler) FindFeedEvents(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var lock sync.Mutex
	res := make([]domain.FeedEvent, 0, limit*2)

	eg.Go(func() error {
		// 查出uid的关注列表
		resp, err := a.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{
			Follower: uid,
			Limit:    10000,
		})
		if err != nil {
			return err
		}

		followeeIds := make([]int64, 0, len(resp.FollowRelations))
		for _, fr := range resp.FollowRelations {
			followeeIds = append(followeeIds, fr.Id)
		}

		evts, err := a.repo.FindPullEventsType(ctx, ArticleEventName, followeeIds, timestamp, limit)
		if err != nil {
			return err
		}

		lock.Lock()
		res = append(res, evts...)
		lock.Unlock()

		return nil
	})

	eg.Go(func() error {
		evts, err := a.repo.FindPushEventsType(ctx, ArticleEventName, uid, timestamp, limit)
		if err != nil {
			return err
		}

		lock.Lock()
		res = append(res, evts...)
		lock.Unlock()

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 查询出了所有数据，需要排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].Ctime.UnixMilli() > res[j].Ctime.UnixMilli()
	})

	minLimit := int(limit)
	if minLimit > len(res) {
		minLimit = len(res)
	}

	return res[:minLimit], nil
}
