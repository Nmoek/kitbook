package events

import (
	"context"
	"github.com/IBM/sarama"
	"kitbook/feed/domain"
	"kitbook/feed/service"
	"kitbook/pkg/logger"
	"kitbook/pkg/saramax"
	"time"
)

type FeedEventConsumer struct {
	svc    service.FeedService
	client sarama.Client
	l      logger.Logger
}

type FeedEvent struct {
	Type     string
	Metadata map[string]string
}

func (f *FeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("feed_event", f.client)
	if err != nil {
		return err
	}
	const topicFeedEvent = "feed_event"
	go func() {
		err2 := cg.Consume(context.Background(), []string{topicFeedEvent}, saramax.NewHandler[FeedEvent](f.Consumer, f.l))
		if err2 != nil {

		}
	}()

	return nil
}

func (f *FeedEventConsumer) Consumer(msg *sarama.ConsumerMessage, event FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return f.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: event.Type,
		Ext:  event.Metadata,
	})

}
