package service

import (
	"context"
	followv1 "kitbook/api/proto/gen/follow/v1"
	"kitbook/feed/domain"
	"kitbook/feed/repository"
	"kitbook/pkg/logger"
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
