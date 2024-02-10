package service

import (
	"context"
	"errors"
	"fmt"
	accountv1 "kitbook/api/proto/gen/account/v1"
	payv1 "kitbook/api/proto/gen/pay/v1"
	"kitbook/pkg/logger"
	"kitbook/reward/domain"
	"kitbook/reward/repository"
	"strconv"
	"strings"
)

type WechatNativeRewardService struct {
	repo   repository.RewardRepository
	payCli payv1.PaymentServiceClient
	acClit accountv1.AccountServiceClient

	l logger.Logger
}

func NewWechatNativeRewardService(repo repository.RewardRepository,
	payCli payv1.PaymentServiceClient,
	acClit accountv1.AccountServiceClient,
	l logger.Logger) RewardService {
	return &WechatNativeRewardService{
		repo:   repo,
		payCli: payCli,
		acClit: acClit,
		l:      l,
	}
}

func (w *WechatNativeRewardService) PreReward(ctx context.Context, rwd domain.Reward) (domain.CodeURL, error) {
	// 是否进行缓存？(可选)
	// 目的: 尽量减少与微信打交道
	codeUrl, err := w.repo.GetCodeUrl(ctx, rwd)
	if err == nil {
		return codeUrl, nil
	}

	w.l.WARN("二维码url缓存未命中",
		logger.Error(err),
		logger.Field{"biz_id", rwd.Target.BizId},
		logger.Field{"biz", rwd.Target.Biz},
		logger.Field{"amt", rwd.Amt})

	// 创建打赏记录
	rwd.Status = domain.RewardStatusInit
	rid, err := w.repo.CreateReward(ctx, rwd)
	if err != nil {
		//TODO: 日志+报错
		return domain.CodeURL{}, err
	}

	// 发起支付请求, 返回支付二维码
	resp, err := w.payCli.PrePay(ctx, &payv1.PrepayRequest{
		// 注意：拼装唯一订单ID需要屏蔽具体传入的业务信息
		BizTradeNo:  fmt.Sprintf("reward-%d", rid),
		Description: fmt.Sprintf("打赏-%s", rwd.Target.BizName),
		Amt: &payv1.Amt{
			Total:    rwd.Amt,
			Currency: "CNY",
		},
		Status: uint32(rwd.Status),
	})
	if err != nil {
		return domain.CodeURL{}, err
	}

	res := domain.CodeURL{
		Rid: rid,
		URL: resp.CodeUrl,
	}

	// 缓存回写
	err = w.repo.SetCodeUrl(ctx, res, rwd)
	if err != nil {
		w.l.ERROR("二维码url缓存回写失败",
			logger.Error(err),
			logger.Int[int64]("biz_id", rwd.Target.BizId),
			logger.Field{"biz", rwd.Target.Biz},
			logger.Field{"amt", rwd.Amt})
	}

	return res, nil
}

func (w *WechatNativeRewardService) UpdateStatus(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {

	rid := w.parseRid(bizTradeNO)
	err := w.repo.UpdateStatus(ctx, rid, status)
	if err != nil {
		return err
	}

	// 完成支付, 进行记账分账流程
	if status == domain.RewardStatusPayed {

		reward, err := w.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}

		// 平台10％的抽成比例
		platformAmt := int64(float64(reward.Amt) * 0.1)

		_, err = w.acClit.Credit(ctx, &accountv1.CreditRequest{
			BizId: rid,
			Biz:   "reward",
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeSystem,
					Amt:         platformAmt,
					Currency:    "CNY",
				},
				{
					Account:     reward.Uid,
					Uid:         reward.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         reward.Amt - platformAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			w.l.ERROR("入账分账失败",
				logger.Error(err),
				logger.Field{"biz_id", rid},
				logger.Field{"biz", "reward"},
				logger.Int[int64]("total_amount", reward.Amt),
				logger.Int[int64]("system_amount", platformAmt))
		}

	}

	return nil
}

// @func: GetReward
// @date: 2024-02-05 23:58:03
// @brief: 查询打赏结果
// @author: Kewin Li
// @receiver w
// @param ctx
// @param rid
// @param uid
func (w *WechatNativeRewardService) GetReward(ctx context.Context, rid int64, uid int64) (domain.Reward, error) {
	// 快路径
	res, err := w.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, nil
	}

	// 不是自己的打赏记录
	if res.Uid != uid {
		return domain.Reward{}, errors.New("非法访问他人的打赏记录")
	}

	//TODO: 如果触发限流、降级就可以不走慢路径

	// 慢路径
	if !res.Complete() {
		pmtResp, err2 := w.payCli.GetPayment(ctx, &payv1.GetPaymentRequest{
			BizTradeNo: w.bizTradeNO(rid),
		})
		if err2 != nil {
			w.l.ERROR("慢路径-主动查询远端支付状态失败",
				logger.Error(err2),
				logger.Int[int64]("rid", rid))
			return res, nil
		}

		// payment --> reward 层级状态转换
		switch pmtResp.Status {
		// 待扫码支付 --> 打赏待支付
		case payv1.Status_PaymentStatusInit:
			res.Status = domain.RewardStatusInit
		// 支付成功 --> 打赏成功
		case payv1.Status_PaymentStatusSuccess:
			res.Status = domain.RewardStatusPayed
		// 支付失败、退款 --> 打赏失败
		case payv1.Status_PaymentStatusFail, payv1.Status_PaymentStatusRefund:
			res.Status = domain.RewardStatusFail
		case payv1.Status_PaymentStatusUnknown:

		}

		err2 = w.UpdateStatus(ctx, w.bizTradeNO(rid), res.Status)
		if err2 != nil {

			w.l.ERROR("慢路径-更新本地状态失败",
				logger.Error(err2),
				logger.Int[int64]("rid", rid))
		}

	}

	return res, nil
}

// @func: parseRid
// @date: 2024-02-08 17:03:37
// @brief: 从唯一支付凭证中去解析获取rid
// @author: Kewin Li
// @receiver w
// @param bizTradeNO
// @return int64
func (w *WechatNativeRewardService) parseRid(bizTradeNO string) int64 {
	idx := strings.Index(bizTradeNO, "-")
	rid, err := strconv.ParseInt(bizTradeNO[idx:], 10, 64)
	if err != nil {
		return -1
	}
	return rid
}

func (w *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}
