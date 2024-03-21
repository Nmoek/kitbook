package service

import (
	"context"
	"kitbook/feed/domain"
)

// FeedService
// @Description: 公有业务逻辑
type FeedService interface {
	CreateFeedEvent(ctx context.Context, event domain.FeedEvent) error
	GetFeedEventList(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
}

// Handler
// @Description: 业务专属逻辑
type Handler interface {
	CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error
	FindFeedEvents(ctx context.Context, uid int64, timestamp int64, limit int64) ([]domain.FeedEvent, error)
}
