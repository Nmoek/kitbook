// Package saramax
// @Description: 同步消费-批量提交
package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"kitbook/pkg/logger"
	"time"
)

// 批量提交阈值
const batchSize = 10

type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	l  logger.Logger
}

func NewBatchHandler[T any](fn func(msgs []*sarama.ConsumerMessage, ts []T) error, l logger.Logger) *BatchHandler[T] {
	return &BatchHandler[T]{
		fn: fn,
		l:  l,
	}
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	msgs := claim.Messages()

	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		timeoutFlag := false
		for i := 0; i < batchSize && timeoutFlag; {
			select {
			case <-ctx.Done():
				//TODO: 会话超时, 日志埋点
				timeoutFlag = true
			case msg, ok := <-msgs:
				// channel被关闭
				if !ok {
					cancel()
					// TODO: 通道关闭错误
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.ERROR("反序列化消息失败",
						logger.Error(err),
						logger.Field{"topic", string(msg.Topic)},
						logger.Int[int32]("partition", msg.Partition),
						logger.Int[int64]("offset", msg.Offset))
					continue
				}

				batch = append(batch, msg)
				ts = append(ts, t)
				// 严格不漏记这里进行递增
				i++
			}

		}

		cancel()

		err := b.fn(batch, ts)
		if err != nil {
			b.l.ERROR("消息业务处理出错",
				logger.Error(err),
				logger.Field{"bizs", ts})

		}

		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}

}
