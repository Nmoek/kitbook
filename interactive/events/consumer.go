package events

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/interactive/repository"
	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"time"
)

const (
	TopicReadEvent = "article_read"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client

	l logger.Logger
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
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
func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewHandler[ReadEvent](i.Consume, i.l))
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
func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return i.repo.IncreaseReadCnt(ctx, "article", event.ArtId)
}

func (i *InteractiveReadEventConsumer) StartV2() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	go func() {

		err2 := cg.Consume(context.Background(), []string{TopicReadEvent}, saramax.NewBatchHandler[ReadEvent](i.BatchConsume, i.l))
		if err2 != nil {
			//TODO: 日志埋点

		}

	}()

	return nil
}

// @func: Consume
// @date: 2023-12-19 03:06:25
// @brief: 帖子模块-实际消费业务处理-批量提交
// @author: Kewin Li
// @receiver i
// @param msg
// @param event
// @return error
func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage, event []ReadEvent) error {
	bizs := make([]string, 0, len(event))
	bizIds := make([]int64, 0, len(event))

	for _, evt := range event {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, evt.ArtId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return i.repo.BatchIncreaseReadCnt(ctx, bizs, bizIds)

}

// ReadEvent
// @Description: 帖子模块-读事件
type ReadEvent struct {
	// 哪一篇文章
	ArtId int64
	// 谁查询的
	UserId int64
}
