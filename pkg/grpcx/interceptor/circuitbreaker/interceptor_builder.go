package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InterceptorBuilder
// @Description: gRPC接入熔断器
type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func NewInterceptorBuilder(breaker circuitbreaker.CircuitBreaker) *InterceptorBuilder {
	return &InterceptorBuilder{
		breaker: breaker,
	}
}

func (i *InterceptorBuilder) BuilderServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		err = i.breaker.Allow()
		if err == nil {
			resp, err = handler(ctx, req)
			// TODO：注意err不代表服务端都是有系统故障, 可能包含业务问题, 需要进一步仔细区分
			if err == nil {
				i.breaker.MarkSuccess()
			} else {
				i.breaker.MarkFailed()
				return nil, status.Errorf(codes.Unavailable, "触发熔断")
			}
		} else {
			i.breaker.MarkFailed()
			return nil, status.Errorf(codes.Unavailable, "触发熔断")
		}

		return
	}
}
