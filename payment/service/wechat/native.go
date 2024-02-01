package wechat

import (
	"context"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"kitbook/payment/domain"
	"kitbook/payment/events"
	"kitbook/payment/repository"
	"kitbook/pkg/logger"
	"time"
)

var errUnknownTransactionState = errors.New("未知微信事务状态")

type NativePaymentService struct {
	appID     string
	mchID     string
	notifyURL string

	repo     repository.PaymentRepository
	svc      *native.NativeApiService
	producer events.Producer

	nativeCBTypeToStatus map[string]domain.PaymentStatus
	l                    logger.Logger
}

func NewNativePaymentService(appID string,
	mchID string,
	repo repository.PaymentRepository,
	svc *native.NativeApiService,
	producer events.Producer,
	l logger.Logger) *NativePaymentService {
	return &NativePaymentService{
		appID:     appID,
		mchID:     mchID,
		notifyURL: "http://xxxx.com/pay/callback/",
		repo:      repo,
		svc:       svc,
		producer:  producer,
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"NOPAY":    domain.PaymentStatusInit,
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFail,
			"REVOKED":  domain.PaymentStatusFail,
			"CLOSED":   domain.PaymentStatusFail,
			"REFUND":   domain.PaymentStatusRefund,
		},
		l: l,
	}
}

// @func: PrePay
// @date: 2024-02-02 02:57:55
// @brief: 发送订单创建请求
// @author: Kewin Li
// @receiver n
// @param ctx
// @param pmt
// @return string
// @return error
func (n *NativePaymentService) PrePay(ctx context.Context, pmt domain.Payment) (string, error) {

	// 创建订单记录
	err := n.repo.CreatePayment(ctx, pmt)
	if err != nil {
		return "", err
	}

	// 调用微信native API
	resp, _, err := n.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(n.appID),
		Mchid:       core.String(n.mchID),
		Description: core.String(pmt.Description),
		OutTradeNo:  core.String(pmt.BizTradeNO),

		TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
		Amount: &native.Amount{
			Total:    core.Int64(pmt.Amt.Total),
			Currency: core.String(pmt.Amt.Currency),
		},
	})

	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, nil
}

// @func: FindExpiredPayment
// @date: 2024-02-02 02:54:50
// @brief: 分批查询过期超时订单
// @author: Kewin Li
// @receiver n
// @param ctx
// @param beforeTime
// @param offset
// @param limit
// @return []domain.Payment
// @return error
func (n *NativePaymentService) FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]domain.Payment, error) {
	return n.repo.FindExpiredPayment(ctx, beforeTime, offset, limit)
}

// @func: HandleCallback
// @date: 2024-02-02 02:54:50
// @brief: 处理支付通知回调
// @author: Kewin Li
// @receiver n
// @param ctx
// @param tnx
// @return error
func (n *NativePaymentService) HandleCallback(ctx context.Context, txn *payments.Transaction) error {
	return n.updateByTxn(ctx, txn)
}

// @func: SyncWechatInfo
// @date: 2024-02-01 23:16:29
// @brief: 主动查询远端订单状态-对账功能
// @author: Kewin Li
// @receiver n
// @param ctx
// @param bizTradeNO
// @return error
func (n *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNO string) error {
	txn, _, err := n.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNO),
	})
	if err != nil {
		return err
	}

	return n.updateByTxn(ctx, txn)

}

// @func: updateByTxn
// @date: 2024-02-01 23:32:14
// @brief: 更新本地信息
// @author: Kewin Li
// @receiver n
// @param ctx
// @param txn
// @return error
func (n *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	status, ok := n.nativeCBTypeToStatus[*txn.TradeState]
	if !ok {
		return fmt.Errorf("%w, 微信状态: %s \n", errUnknownTransactionState, txn.TradeState)
	}

	err := n.repo.UpdatePayment(ctx, domain.Payment{
		TxnID:      *txn.TransactionId,
		BizTradeNO: *txn.OutTradeNo,
		Status:     status,
	})
	if err != nil {
		return err
	}

	// 消息队列 向打赏模块发送消息
	err = n.producer.ProducePaymentEvent(ctx, events.PaymentEvent{
		BizTradeNO: *txn.OutTradeNo,
		Status:     status.AstoUint8(),
	})
	//TODO: 消息发送失败怎么办？ 部分失败问题
	// 至少发送成功一次。重复消息怎么处理？

	if err != nil {
		n.l.ERROR("支付事件消息发送失败",
			logger.Error(err),
			logger.Field{"biz_trade_no", *txn.OutTradeNo},
			logger.Field{"txn_status", *txn.TradeState},
			logger.Int[uint8]("status", status.AstoUint8()))
	}

	return nil
}
