package repository

import (
	"context"
	"encoding/json"
	"kitbook/feed/domain"
	"kitbook/feed/repository/dao"
	"time"
)

type FeedEventRepository interface {
	CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error
	CreatePullEvent(ctx context.Context, event domain.FeedEvent) error

	// 批量获取拉事件, 关注的人的发件箱中数据
	FindPullEvents(ctx context.Context, ids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
	FindPushEvents(ctx context.Context, id int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)

	FindPushEventsType(ctx context.Context, typ string, id int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
	FindPullEventsType(ctx context.Context, typ string, ids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
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
// @brief: 创建批量推事件
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
// @brief: 创建单个拉事件
// @author: Kewin Li
// @receiver f
// @param ctx
// @param events
// @return error
func (f *feedEventRepository) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return f.pullDao.CreatePullEvent(ctx, f.ConvertsPullEventDao(&event))
}

func (f *feedEventRepository) FindPushEvents(ctx context.Context, id int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	resDao, err := f.pushDao.FindPushEvents(ctx, id, timestamp, limit)

	if err != nil {
		return nil, err
	}
	return f.ConvertsPushEventsDomain(resDao), nil
}

func (f *feedEventRepository) FindPullEvents(ctx context.Context, ids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	resDao, err := f.pullDao.FindPullEvents(ctx, ids, timestamp, limit)
	if err != nil {
		return nil, err
	}

	return f.ConvertsPullEventsDomain(resDao), nil
}

func (f *feedEventRepository) FindPushEventsType(ctx context.Context, typ string, id int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	resDao, err := f.pushDao.FindPushEventsType(ctx, typ, id, timestamp, limit)
	if err != nil {
		return nil, err
	}

	return f.ConvertsPushEventsDomain(resDao), nil
}

func (f *feedEventRepository) FindPullEventsType(ctx context.Context, typ string, ids []int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	resDao, err := f.pullDao.FindPullEventsType(ctx, typ, ids, timestamp, limit)
	if err != nil {
		return nil, err
	}

	return f.ConvertsPullEventsDomain(resDao), nil
}

func (f *feedEventRepository) ConvertsPushEventDao(event *domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
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

func (f *feedEventRepository) ConvertsPullEventDao(event *domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Uid:     event.Uid,
		Content: string(val),
		Ctime:   event.Ctime.UnixMilli(),
	}
}

func (f *feedEventRepository) ConvertsPullEventDomain(event *dao.FeedPullEvent) domain.FeedEvent {
	var ext domain.ExtendFields
	_ := json.Unmarshal([]byte(event.Content), &ext)

	return domain.FeedEvent{
		Id:    event.Id,
		Uid:   event.Uid,
		Type:  event.Type,
		Ext:   ext,
		Ctime: time.UnixMilli(event.Ctime),
	}
}

func (f *feedEventRepository) ConvertsPullEventsDomain(events []dao.FeedPullEvent) []domain.FeedEvent {
	res := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		res = append(res, f.ConvertsPullEventDomain(&e))
	}

	return res
}

func (f *feedEventRepository) ConvertsPushEventDomain(event *dao.FeedPushEvent) domain.FeedEvent {
	var ext domain.ExtendFields
	_ := json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		Id:    event.Id,
		Type:  event.Type,
		Uid:   event.Uid,
		Ext:   ext,
		Ctime: time.UnixMilli(event.Ctime),
	}
}

func (f *feedEventRepository) ConvertsPushEventsDomain(events []dao.FeedPushEvent) []domain.FeedEvent {
	res := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		res = append(res, f.ConvertsPushEventDomain(&e))
	}

	return res
}
