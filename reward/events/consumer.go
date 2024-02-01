package events

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/reward/domain"
	"kitbook/reward/service"

	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"strings"
	"time"
)

const (
	TopicPaymentEvent = "payment_events"
)

type PaymentEventConsumer struct {
	client sarama.Client
	svc    service.RewardService
	l      logger.Logger
}

func NewPaymentEventConsumer(svc service.RewardService,
	client sarama.Client,
	l logger.Logger) *PaymentEventConsumer {
	return &PaymentEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

// @func: Start
// @date: 2023-12-17 20:25:40
// @brief: 启动消费
// @author: Kewin Li
// @receiver i
// @return error
func (p *PaymentEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("reward", p.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicPaymentEvent}, saramax.NewHandler[PaymentEvent](p.Consume, p.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: Consume
// @date: 2024-02-05 23:19:21
// @brief: 打赏模块-实际消费业务处理-支付成功通知
// @author: Kewin Li
// @receiver p
// @param msg
// @param event
// @return error
func (p *PaymentEventConsumer) Consume(msg *sarama.ConsumerMessage, event PaymentEvent) error {
	if !strings.HasPrefix(event.BizTradeNO, "reward") {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	return p.svc.UpdateStatus(ctx, event.BizTradeNO, event.ToDomainStatus())

}

func (p *PaymentEventConsumer) StartV2() error {
	cg, err := sarama.NewConsumerGroupFromClient("reward", p.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicPaymentEvent}, saramax.NewBatchHandler[PaymentEvent](p.BatchConsume, p.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: BatchConsume
// @date: 2024-02-05 23:19:48
// @brief: 打赏模块-实际消费业务处理-支付成功通知批量提交
// @author: Kewin Li
// @receiver p
// @param msgs
// @param event
// @return error
func (p *PaymentEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage, events []PaymentEvent) error {
	//bizs := make([]string, 0, len(event))
	//bizIds := make([]int64, 0, len(event))
	//
	//for _, evt := range event {
	//	bizs = append(bizs, "article")
	//	bizIds = append(bizIds, evt.ArtId)
	//}
	//
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//
	//return p.repo.BatchIncreaseReadCnt(ctx, bizs, bizIds)
	panic("todo")
}

// PaymentEvent
// @Description: 支付成功事件
type PaymentEvent struct {
	BizTradeNO string
	Status     uint8
}

func (p PaymentEvent) ToDomainStatus() domain.RewardStatus {
	switch p.Status {
	// PaymentStatusInit 支付初始化
	case 1:
		return domain.RewardStatusInit
	// PaymentStatusSuccess 支付成功
	case 2:
		return domain.RewardStatusPayed
	// PaymentStatusFail、PaymentStatusRefund 支付失败/退款
	case 3, 4:
		return domain.RewardStatusFail
	default:
		return domain.RewardStatusUnknown
	}
}
