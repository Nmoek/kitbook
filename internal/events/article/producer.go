// Package article
// @Description: 领域事件-帖子模块消息发送
package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	TopicReadEvent = "article_read"
)

type Producer interface {
	ProducerReadEvent(event ReadEvent) error
}

// SaramaSyncProducer
// @Description: 帖子模块-消息同步发送
type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

// @func: ProducerReadEvent
// @date: 2023-12-17 19:56:31
// @brief: 帖子模块读事件-阅读数+1消息
// @author: Kewin Li
// @receiver s
func (s *SaramaSyncProducer) ProducerReadEvent(event ReadEvent) error {
	val, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})

	// TODO: 日志埋点, 向哪个分区、哪一段偏移发送了消息

	return err
}

// ReadEvent
// @Description: 帖子模块-读事件
type ReadEvent struct {
	// 哪一篇文章
	ArtId int64
	// 谁查询的
	UserId int64
}
