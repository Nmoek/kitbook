package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed sliding_window.lua
var luaSlidingWindow string

// RedisSlidingWindowLimiter
// @Description: 滑动窗口限流器
type RedisSlidingWindowLimiter struct {
	cmd        redis.Cmdable
	windowSize time.Duration
	threshold  int
}

func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, windowSize time.Duration, threshold int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:        cmd,
		windowSize: windowSize,
		threshold:  threshold,
	}
}

func (r *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(
		ctx,
		luaSlidingWindow,
		[]string{key},
		r.windowSize.Milliseconds(),
		r.threshold,
		time.Now().UnixMilli()).Bool()
}
