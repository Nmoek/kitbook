package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "kitbook/feed/grpc"
	"kitbook/pkg/grpcx"
	"kitbook/pkg/logger"
)

func InitGRpcServer(feedSvc *grpc2.FeedServiceServer,
	l logger.Logger) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcd_addr"`
		Port     int    `yaml:"port"`
		BizName  string `yaml:"biz_name"`
	}
	var cfg Config

	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	feedSvc.Register(s)
	return &grpcx.Server{
		Server:   s,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		BizName:  "feed",
		L:        l,
	}
}
