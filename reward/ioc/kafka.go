package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"kitbook/internal/events"
	events2 "kitbook/reward/events"
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

// 注意： wire没有办法找到所有同类实现
func InitConsumers(c *events2.PaymentEventConsumer) []events.Consumer {

	return []events.Consumer{c}

}
