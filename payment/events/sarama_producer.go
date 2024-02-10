package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const (
	TopicPaymentEvent = "payment_events"
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) *SaramaSyncProducer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

func (s *SaramaSyncProducer) ProducePaymentEvent(ctx context.Context, event PaymentEvent) error {
	val, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicPaymentEvent,
		Value: sarama.StringEncoder(val),
	})

	return err
}
