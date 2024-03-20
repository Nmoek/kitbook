package service

import (
	"context"
	"kitbook/follow/domain"
	"kitbook/follow/repository"
	"kitbook/pkg/logger"
)

type ArticleFollowService struct {
	repo repository.FollowRepository

	l logger.Logger
}

func NewArticleFollowService(repo repository.FollowRepository, l logger.Logger) FollowService {
	return &ArticleFollowService{
		repo: repo,
		l:    l,
	}
}

func (a *ArticleFollowService) Follow(ctx context.Context, followee, follower int64) error {
	return a.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}

func (a *ArticleFollowService) CancelFollow(ctx context.Context, followee, follower int64) error {
	return a.repo.InactiveFollowRelation(ctx, followee, follower)
}

func (a *ArticleFollowService) GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error) {
	return a.repo.GetFollowee(ctx, follower, offset, limit)
}

func (a *ArticleFollowService) GetFollower(ctx context.Context, followee int64, offset int64, limit int64) ([]domain.FollowRelation, error) {
	return a.repo.GetFollowee(ctx, followee, offset, limit)
}

func (a *ArticleFollowService) FollowInfo(ctx context.Context, followee, follower int64) (domain.FollowRelation, error) {
	return a.repo.FollowInfo(ctx, followee, follower)
}

func (a *ArticleFollowService) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return a.repo.GetFollowStatics(ctx, uid)
}
