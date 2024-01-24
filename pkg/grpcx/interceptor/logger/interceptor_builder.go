package logger

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kitbook/pkg/grpcx/interceptor"
	"kitbook/pkg/logger"
	"runtime"
	"time"
)

type InterceptorBuilder struct {
	l logger.Logger
	interceptor.Builder
}

func (i *InterceptorBuilder) BuildServerUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		start := time.Now()
		event := "normal"
		defer func() {
			// 最终输出日志
			cost := time.Since(start)

			if rec := recover(); rec != nil {
				switch res := rec.(type) {
				case error:
					err = res
				default:
					err = fmt.Errorf("%v", rec)
				}
				// 获取panic堆栈
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				err = status.New(codes.Internal, "panic err: "+err.Error()).Err()

				event = "recover"
			}

			fields := []logger.Field{
				// unary调用一次  stream流式调用
				logger.Field{"type", "unary"},
				logger.Int[int64]("cost_time", cost.Milliseconds()),
				logger.Field{"event", event},
				logger.Field{"method", info.FullMethod},
				logger.Field{"perr", i.PeerName(ctx)},
				logger.Field{"peer_ip", i.PeerIP(ctx)},
			}

			st, _ := status.FromError(err)
			if st != nil {
				fields = append(fields, logger.Field{"code", st.Code().String()})
				fields = append(fields, logger.Field{"cods_msg", st.Message()})

			}

			i.l.INFO("RPC服务调用", fields...)

		}()
		resp, err = handler(ctx, req)
		return

	}
}
