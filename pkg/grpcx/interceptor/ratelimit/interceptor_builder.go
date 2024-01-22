package ratelimit

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	intrv1 "kitbook/api/proto/gen/intr/v1"
	"kitbook/pkg/limiter"
	"strings"
)

// InterceptorBuilder
// @Description: gRPC接入限流器
type InterceptorBuilder struct {
	limiter limiter.Limiter
	key     string
}

func NewInterceptorBuilder(limiter limiter.Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{
		limiter: limiter,
		key:     key,
	}
}

// @func: BuildServerUnaryInterceptor
// @date: 2024-01-21 22:01:37
// @brief: 应用级别限流
// @author: Kewin Li
// @receiver i
// @return grpc.UnaryServerInterceptor
func (i *InterceptorBuilder) BuildServerUnaryInterceptorApp() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		limited, err := i.limiter.Limit(ctx, i.key)
		if err != nil {
			// TODO: 思考更多可能性
			// 保守做法: 返回错误
			// 激进做法: 继续执行不限流
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}

		// 触发限流
		if limited {
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}

		return handler(ctx, req)

	}
}

// @func: BuildServerUnaryInterceptorService
// @date: 2024-01-21 22:02:01
// @brief: 服务级别限流
// @author: Kewin Li
// @receiver i
// @return grpc.UnaryServerInterceptor
func (i *InterceptorBuilder) BuildServerUnaryInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		if strings.Contains(info.FullMethod, "/Uservice") {

			limited, err := i.limiter.Limit(ctx, i.key)
			if err != nil {
				// TODO: 思考更多可能性
				// 保守做法: 返回错误
				// 激进做法: 继续执行不限流
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}

			// 触发限流
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}
		return handler(ctx, req)

	}
}

// @func: BuildServerUnaryInterceptorBusiness
// @date: 2024-01-21 22:07:08
// @brief: 业务级别限流（演示）
// @author: Kewin Li
// @receiver i
// @return grpc.UnaryServerInterceptor
func (i *InterceptorBuilder) BuildServerUnaryInterceptorBusiness() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		if getReq, ok := req.(*intrv1.LikeRequest); ok {
			key := fmt.Sprintf("limiter:user:like:%d", getReq.BizId)

			limited, err := i.limiter.Limit(ctx, key)
			if err != nil {
				// TODO: 思考更多可能性
				// 保守做法: 返回错误
				// 激进做法: 继续执行不限流
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}

			// 触发限流
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}
		return handler(ctx, req)

	}
}
