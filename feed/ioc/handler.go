package ioc

import (
	followv1 "kitbook/api/proto/gen/follow/v1"
	"kitbook/feed/repository"
	"kitbook/feed/service"
	"kitbook/pkg/logger"
)

func InitHandler(repo repository.FeedEventRepository, cli followv1.FollowServiceClient, l logger.Logger) map[string]service.Handler {
	artHandler := service.NewArticleHandler(repo, cli, l)
	likeHandler := service.NewLikeHandler(repo, l)
	followHandler := service.NewFollowHandler(repo, l)

	return map[string]service.Handler{
		service.ArticleEventName: artHandler,
		service.FollowEventName:  followHandler,
		service.LikeEventName:    likeHandler,
	}
}
