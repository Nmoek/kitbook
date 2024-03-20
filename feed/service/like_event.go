package service

import (
	"context"
	"kitbook/feed/domain"
	"kitbook/feed/repository"
	"kitbook/pkg/logger"
	"time"
)

const (
	LikeEventName = "like_event"
)

type LikeHandler struct {
	repo repository.FeedEventRepository
	l2   logger.Logger
}

func NewLikeHandler(repo repository.FeedEventRepository, l2 logger.Logger) *LikeHandler {
	return &LikeHandler{
		repo: repo,
		l2:   l2,
	}
}

func (l *LikeHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	// 字段检验(可做可不做)

	// 被点赞人的Id
	uid, _ := ext.Get("liked").AsInt64()

	return l.repo.CreatePushEvents(ctx, []domain.FeedEvent{
		{
			Type:  LikeEventName,
			Uid:   uid,
			Ext:   ext,
			Ctime: time.Now(),
		},
	})
}
