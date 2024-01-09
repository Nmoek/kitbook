package ioc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	grpc2 "kitbook/interactive/grpc"
	"kitbook/pkg/grpcx"
)

func InitGRpcServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	return &grpcx.Server{
		Server: s,
		Addr:   viper.GetString("grpc.server.addr"),
	}
}
