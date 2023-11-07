package ioc

import (
	"github.com/redis/go-redis/v9"
	"kitbook/pkg/limiter"
	"time"
)

func InitLimiter(client redis.Cmdable) limiter.Limiter {
	return limiter.NewRedisSlidingWindowLimiter(client, time.Second, 1000)
}
