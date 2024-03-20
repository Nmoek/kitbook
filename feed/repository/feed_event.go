package repository

import (
	"context"
	"encoding/json"
	"kitbook/feed/domain"
	"kitbook/feed/repository/dao"
)

type FeedEventRepository interface {
	CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error
	CreatePullEvent(ctx context.Context, event domain.FeedEvent) error
}

type feedEventRepository struct {
	pullDao dao.FeedPullEventDao
	pushDao dao.FeedPushEventDao
}

func NewFeedEventRepository(pullDao dao.FeedPullEventDao, pushDao dao.FeedPushEventDao) FeedEventRepository {
	return &feedEventRepository{
		pullDao: pullDao,
		pushDao: pushDao,
	}
}

// @func: CreatePushEvents
// @date: 2024-03-20 22:15:14
// @brief: 批量推事件
// @author: Kewin Li
// @receiver f
// @param ctx
// @param events
// @return error
func (f *feedEventRepository) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	return f.pushDao.CreatePushEvents(ctx, f.ConvertsPushEventsDao(events))
}

// @func: CreatePullEvent
// @date: 2024-03-20 23:02:55
// @brief: 单个拉事件
// @author: Kewin Li
// @receiver f
// @param ctx
// @param events
// @return error
func (f *feedEventRepository) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, f.ConvertsPullEventDao(&event))

}

func (f *feedEventRepository) ConvertsPushEventDao(event *domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
		Uid:     event.Uid,
		Content: string(val),
		Ctime:   event.Ctime.UnixMilli(),
	}
}

func (f *feedEventRepository) ConvertsPullEventDao(event *domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Uid:     event.Uid,
		Content: string(val),
		Ctime:   event.Ctime.UnixMilli(),
	}
}

func (f *feedEventRepository) ConvertsPushEventsDao(events []domain.FeedEvent) []dao.FeedPushEvent {
	res := make([]dao.FeedPushEvent, 0, len(events))
	for _, t := range events {
		res = append(res, f.ConvertsPushEventDao(&t))
	}

	return res
}
