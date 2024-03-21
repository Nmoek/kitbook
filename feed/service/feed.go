package service

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	followv1 "kitbook/api/proto/gen/follow/v1"
	"kitbook/feed/domain"
	"kitbook/feed/repository"
	"kitbook/pkg/logger"
	"sort"
	"sync"
)

type feedService struct {
	followClient followv1.FollowServiceClient
	repo         repository.FeedEventRepository
	handlers     map[string]Handler
	l            logger.Logger
}

func NewFeedService(handlers map[string]Handler, l logger.Logger) FeedService {
	return &feedService{
		handlers: handlers,
		l:        l,
	}
}

func (f *feedService) Register(key string, h Handler) {
	f.handlers[key] = h
}

func (f *feedService) CreateFeedEvent(ctx context.Context, event domain.FeedEvent) error {
	handler, ok := f.handlers[event.Type]
	if !ok {
		// TODO: 可以考虑兜底机制, 设置一个默认Handler
		return fmt.Errorf("未能找到对应Handler, %s", event.Type)
	}

	return handler.CreateFeedEvent(ctx, event.Ext)

}

// @func: GetFeedEventList
// @date: 2024-03-22 01:20:47
// @brief: 查询——Service层直接查询
// @author: Kewin Li
// @receiver f
// @param ctx
// @param uid
// @param timestamp
// @param limit
// @return []domain.FeedEvent
// @return error
func (f *feedService) GetFeedEventListV1(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var lock sync.Mutex
	res := make([]domain.FeedEvent, 0, limit*2)

	eg.Go(func() error {
		// 查出uid的关注列表
		resp, err := f.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{
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

		evts, err := f.repo.FindPullEvents(ctx, followeeIds, timestamp, limit)
		if err != nil {
			return err
		}

		lock.Lock()
		res = append(res, evts...)
		lock.Unlock()

		return nil
	})

	eg.Go(func() error {
		evts, err := f.repo.FindPushEvents(ctx, uid, timestamp, limit)
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

// @func: GetFeedEventList
// @date: 2024-03-22 02:01:39
// @brief: 查询——Handler上查询
// @author: Kewin Li
// @receiver f
// @param ctx
// @param uid
// @param timestamp
// @param limit
// @return []domain.FeedEvent
// @return error
func (f *feedService) GetFeedEventList(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var lock sync.Mutex
	res := make([]domain.FeedEvent, 0, int(limit)*len(f.handlers))

	for _, h := range f.handlers {
		eg.Go(func() error {
			evts, err := h.FindFeedEvents(ctx, uid, timestamp, limit)
			if err != nil {
				return err
			}

			lock.Lock()
			res = append(res, evts...)
			lock.Unlock()

			return nil
		})

	}
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
