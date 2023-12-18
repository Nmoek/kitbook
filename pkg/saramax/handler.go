// Package saramax
// @Description: 对sarama消费API二次封装
package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"kitbook/pkg/logger"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error

	l logger.Logger
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error, l logger.Logger) *Handler[T] {
	return &Handler[T]{
		fn: fn,
		l:  l,
	}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil

}

// @func: ConsumeClaim
// @date: 2023-12-17 20:21:01
// @brief: 消费业务处理
// @author: Kewin Li
// @receiver h
// @param session
// @param claim
// @return error
func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()

	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.ERROR("反序列化消息失败",
				logger.Error(err),
				logger.Field{"topic", string(msg.Topic)},
				logger.Int[int32]("partition", msg.Partition),
				logger.Int[int64]("offset", msg.Offset))
		}

		err = h.fn(msg, t)
		if err != nil {
			h.l.ERROR("消息业务处理出错",
				logger.Error(err),
				logger.Field{"biz", t},
				logger.Field{"topic", msg.Topic},
				logger.Int[int32]("partition", msg.Partition),
				logger.Int[int64]("offset", msg.Offset))
		}

		session.MarkMessage(msg, "")
	}

	return nil
}
