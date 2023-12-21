// Package article
// @Description: 浏览历史记录消费业务
package article

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"time"
)

type HistoryRecordConsumer struct {
	repo   repository.HistoryRepository
	client sarama.Client
	l      logger.Logger
}

func NewHistoryRecordConsumer(repo repository.InteractiveRepository,
	client sarama.Client,
	l logger.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		repo:   repo,
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
func (h *HistoryRecordConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", h.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](h.Consume, h.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: Consume
// @date: 2023-12-17 20:31:03
// @brief: 帖子模块-实际消费业务处理-阅读数+1
// @author: Kewin Li
// @receiver i
// @param msg
// @param events
// @return error
func (h *HistoryRecordConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return h.repo.AddRecord(ctx, domain.HistoryRecord{
		BizId:  event.ArtId,
		Biz:    "article", // 帖子业务标识
		UserId: event.UserId,
	})
}

func (h *HistoryRecordConsumer) StartV1() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", h.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewBatchHandler[ReadEvent](h.BatchConsume, h.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: BatchConsume
// @date: 2023-12-19 12:46:41
// @brief: 帖子模块-实际消费业务处理-批量提交
// @author: Kewin Li
// @receiver h
// @param msgs
// @param event
// @return error
func (h *HistoryRecordConsumer) BatchConsume(msgs []*sarama.ConsumerMessage, event []ReadEvent) error {
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
	//return h.repo.BatchAddRecord(ctx, bizs, bizIds)
	panic("接口位实现")
}
