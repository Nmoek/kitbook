package startup

import (
	"github.com/IBM/sarama"
	events2 "kitbook/interactive/events"
	"kitbook/internal/events"
)

func InitSaramaClient() sarama.Client {

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9094"}, cfg)

	if err != nil {
		panic(err)
	}

	return client

}

func InitSyncProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

// 注意： wire没有办法找到所有同类实现
func InitConsumers(c *events2.InteractiveReadEventConsumer) []events.Consumer {

	return []events.Consumer{c}

}
