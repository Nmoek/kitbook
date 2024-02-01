package grpc

import (
	"google.golang.org/grpc"
	payv1 "kitbook/api/proto/gen/pay/v1"
	"kitbook/payment/service/wechat"
)

type PaymentServiceServer struct {
	payv1.UnimplementedPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewPaymentServiceServer(svc *wechat.NativePaymentService) *PaymentServiceServer {
	return &PaymentServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (p *PaymentServiceServer) Register(s *grpc.Server) {
	payv1.RegisterPaymentServiceServer(s, p)
}
