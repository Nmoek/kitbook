package ioc

import (
	olivere "github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

func InitES() *olivere.Client {
	type Config struct {
		Addr string `json:"addr"`
	}
	var cfg Config

	err := viper.UnmarshalKey("es.server.addr", &cfg)
	if err != nil {
		panic(err)
	}

	client, err := olivere.NewClient(olivere.SetURL(cfg.Addr),
		olivere.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return client
}
