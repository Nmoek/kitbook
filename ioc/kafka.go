package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"kitbook/internal/events"
	"kitbook/internal/events/article"
)

func InitSaramaClient() sarama.Client {
	// 配置管理
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}

	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)

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
func InitConsumers(c *article.InteractiveReadEventConsumer) []events.Consumer {

	return []events.Consumer{c}

}
