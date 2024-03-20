package service

import (
	"context"
	"fmt"
	"kitbook/feed/domain"
	"kitbook/pkg/logger"
)

type feedService struct {
	handlers map[string]Handler
	l        logger.Logger
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
