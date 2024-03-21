package service

import (
	"context"
	"kitbook/feed/domain"
	"kitbook/feed/repository"
	"kitbook/pkg/logger"
	"time"
)

const (
	FollowEventName = "followee_event"
)

type FollowHandler struct {
	repo repository.FeedEventRepository
	l    logger.Logger
}

func NewFollowHandler(repo repository.FeedEventRepository, l logger.Logger) Handler {
	return &FollowHandler{
		repo: repo,
		l:    l,
	}
}

func (f *FollowHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	followee, _ := ext.Get("followee").AsInt64()

	return f.repo.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			Type:  FollowEventName,
			Uid:   followee,
			Ctime: time.Now(),
			Ext:   ext,
		},
	})
}

func (f *FollowHandler) FindFeedEvents(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error) {
	return f.repo.FindPushEventsType(ctx, FollowEventName, uid, timestamp, limit)
}
