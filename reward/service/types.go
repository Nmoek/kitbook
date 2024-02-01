package service

import (
	"context"
	"kitbook/reward/domain"
)

type RewardService interface {
	PreReward(ctx context.Context, pmt domain.Reward) (domain.CodeURL, error)
	UpdateStatus(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error
	GetReward(ctx context.Context, rid int64, uid int64) (domain.Reward, error)
}
