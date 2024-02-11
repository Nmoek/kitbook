package cache

import "github.com/redis/go-redis/v9"

type PaymentCache interface {
}

type RedisPaymentCache struct {
	client redis.Cmdable
}

func NewRedisPaymentCache(client redis.Cmdable) PaymentCache {
	return &RedisPaymentCache{
		client: client,
	}
}
