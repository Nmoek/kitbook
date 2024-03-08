package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "kitbook/comment/grpc"
	"kitbook/pkg/grpcx"
	"kitbook/pkg/logger"
)

func InitGRpcServer(commentSvc *grpc2.ArticleCommentServiceServer, l logger.Logger) *grpcx.Server {
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
	commentSvc.Register(s)
	return &grpcx.Server{
		Server:   s,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		BizName:  cfg.BizName,
		L:        l,
	}
}
