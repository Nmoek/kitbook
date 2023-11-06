// Package ratelimit
// @Description: 基于gin的限流插件
package ratelimit

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

//go:embed sliding_window.lua
var luaSlidingWindow string

type MiddlewareBuilder struct {
	client     redis.Cmdable
	keyGenFunc func(ctx *gin.Context) string
	windowSize time.Duration
	threshold  int
}

func NewMiddlewareBuilder(client redis.Cmdable, windowSize time.Duration, threshold int) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		client:     client,
		windowSize: windowSize,
		threshold:  threshold,
		keyGenFunc: func(ctx *gin.Context) string {
			return "ip-limiter:" + ctx.ClientIP()
		},
	}
}

// 业务代码中实现

//func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		// 当前时间戳
//		key := m.keyGenFunc(ctx)
//		now := time.Now().UnixNano()
//		windowStart := fmt.Sprintf("%d", now-m.windowSize.Nanoseconds())
//
//		//TODO: 1个用户请求在这里分裂成至多4个redis请求，存在redis性能问题
//
//		//清理不在窗口范围的时间戳
//		err := m.client.ZRemRangeByScore(ctx, key, "0", windowStart).Err()
//		if err != nil {
//			ctx.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//
//		// 统计窗口中还剩多少请求
//		reqs, err := m.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
//			Min: windowStart,
//			Max: fmt.Sprintf("%d", now),
//		}).Result()
//
//		// 超过阈值被限流
//		// TODO: 存在业务并发问题(使用lua脚本解决)
//		if len(reqs) >= m.threshold {
//			ctx.AbortWithStatus(http.StatusTooManyRequests)
//			return
//		}
//
//		err = m.client.ZAddNX(ctx, key, redis.Z{
//			Score:  float64(now),
//			Member: "", // 随便放一个值
//		}).Err()
//
//		if err != nil {
//			return
//		}
//
//		// 添加过期时间
//		m.client.Expire(ctx, key, m.windowSize)
//		// TODO: 日志埋点
//
//		ctx.Next()
//
//	}
//}

func (m MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := m.limit(ctx)
		if err != nil {
			// redis崩溃了是否应该限流
			// 保守策略，返回错误，不接收请求
			ctx.AbortWithStatus(http.StatusInternalServerError)

			// 1. 尽可能保证可用性，继续处理请求，不限流
			// 2. 集群限流失效，启用单实例/单机限流
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
	}
}

func (m MiddlewareBuilder) limit(ctx *gin.Context) (bool, error) {
	key := m.keyGenFunc(ctx)
	return m.client.Eval(ctx,
		luaSlidingWindow,
		[]string{key},
		m.windowSize.Milliseconds(),
		m.threshold,
		time.Now().UnixMilli()).Bool()
}
