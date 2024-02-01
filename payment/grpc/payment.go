package grpc

import (
	"context"
	"google.golang.org/grpc"
	payv1 "kitbook/api/proto/gen/pay/v1"
	"kitbook/payment/domain"
	"kitbook/payment/service/wechat"
)

// NativePaymentServiceServer
// @Description: 微信native API接口远程服务
type NativePaymentServiceServer struct {
	payv1.UnimplementedPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewPaymentServiceServer(svc *wechat.NativePaymentService) *NativePaymentServiceServer {
	return &NativePaymentServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (p *NativePaymentServiceServer) Register(s *grpc.Server) {
	payv1.RegisterPaymentServiceServer(s, p)
}

// @func: PrePay
// @date: 2024-02-05 00:54:18
// @brief: 远程调用-向微信发起支付请求
// @author: Kewin Li
// @receiver p
// @param ctx
// @param request
// @return *payv1.PrepayResponse
// @return error
func (p *NativePaymentServiceServer) PrePay(ctx context.Context, request *payv1.PrepayRequest) (*payv1.PrepayResponse, error) {
	codeUrl, err := p.svc.PrePay(ctx, domain.Payment{
		BizTradeNO: request.BizTradeNo,
		Amt: domain.Amount{
			Total:    request.Amt.Total,
			Currency: request.Amt.Currency,
		},
		Description: request.Description,
	})

	return &payv1.PrepayResponse{
		CodeUrl: codeUrl,
	}, err

}
