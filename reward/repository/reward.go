package repository

import (
	"context"
	"kitbook/reward/domain"
	"kitbook/reward/repository/cache"
	"kitbook/reward/repository/dao"
)

type RewardRepository interface {
	CreateReward(ctx context.Context, rwd domain.Reward) (int64, error)
	UpdateReward(ctx context.Context, rwd domain.Reward) error
	UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error

	SetCodeUrl(ctx context.Context, codeUrl domain.CodeURL, rwd domain.Reward) error
	GetCodeUrl(ctx context.Context, rwd domain.Reward) (domain.CodeURL, error)
	DelCodeUrl(ctx context.Context, rwd domain.Reward) error
	GetReward(ctx context.Context, rid int64) (domain.Reward, error)
}

type WechatNativeRewardRepository struct {
	dao   dao.RewardDao
	cache cache.RewardCache
}

func NewWechatNativeRewardRepository(dao dao.RewardDao, cache cache.RewardCache) RewardRepository {
	return &WechatNativeRewardRepository{
		dao:   dao,
		cache: cache,
	}
}

func (w *WechatNativeRewardRepository) CreateReward(ctx context.Context, rwd domain.Reward) (int64, error) {
	return w.dao.Insert(ctx, w.ConvertsDaoReward(&rwd))
}

func (w *WechatNativeRewardRepository) UpdateReward(ctx context.Context, rwd domain.Reward) error {
	return w.dao.Update(ctx, w.ConvertsDaoReward(&rwd))
}
func (w *WechatNativeRewardRepository) UpdateStatus(ctx context.Context, rid int64, status domain.RewardStatus) error {
	// 删除缓存
	//err := w.cache.DelCodeUrl(ctx, bizId, biz)
	//if err != nil {
	//	//TODO: 缓存删除失败
	//}

	return w.dao.UpdateStatus(ctx, rid, status.AsToUint8())

}

func (w *WechatNativeRewardRepository) GetReward(ctx context.Context, rid int64) (domain.Reward, error) {

	res, err := w.dao.GetByRid(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	return w.ConvertsDomainReward(&res), nil
}

func (w *WechatNativeRewardRepository) SetCodeUrl(ctx context.Context, codeUrl domain.CodeURL, rwd domain.Reward) error {
	return w.cache.SetCodeUrl(ctx, codeUrl, rwd)
}

func (w *WechatNativeRewardRepository) GetCodeUrl(ctx context.Context, rwd domain.Reward) (domain.CodeURL, error) {
	return w.cache.GetCodeUrl(ctx, rwd)
}

func (w *WechatNativeRewardRepository) DelCodeUrl(ctx context.Context, rwd domain.Reward) error {
	return w.cache.DelCodeUrl(ctx, rwd.Target.BizId, rwd.Target.Biz)
}

func (w *WechatNativeRewardRepository) ConvertsDomainReward(rwd *dao.Reward) domain.Reward {
	return domain.Reward{
		Uid: rwd.Uid,
		Target: domain.Target{
			BizId:   rwd.BizId,
			BizName: rwd.BizName,
			Biz:     rwd.Biz,
		},
		Status: domain.RewardStatus(rwd.Status),
		Amt:    rwd.Amt,
	}
}

func (w *WechatNativeRewardRepository) ConvertsDaoReward(rwd *domain.Reward) dao.Reward {
	return dao.Reward{
		BizId:     rwd.Target.BizId,
		Biz:       rwd.Target.Biz,
		BizName:   rwd.Target.BizName,
		TargetUid: rwd.Target.Uid,
		Status:    rwd.Status.AsToUint8(),
		Uid:       rwd.Uid,
		Amt:       rwd.Amt,
	}
}
