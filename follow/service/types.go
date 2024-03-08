package service

import (
	"context"
	"kitbook/follow/domain"
)

type FollowService interface {
	Follow(ctx context.Context, followee, follower int64) error
	CancelFollow(ctx context.Context, followee, follower int64) error

	GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error)
	FollowInfo(ctx context.Context, followee, follower int64) (domain.FollowRelation, error)
}
