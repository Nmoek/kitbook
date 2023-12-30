package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"kitbook/internal/domain"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
}

type RedisRankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRedisRankingCache(client redis.Cmdable) RankingCache {
	return &RedisRankingCache{
		client:     client,
		key:        "ranking:top_n",
		expiration: 3 * time.Minute, // 3min 缓存过期时间
	}
}

// @func: Set
// @date: 2023-12-29 23:01:36
// @brief: 热榜缓存-插入数据
// @author: Kewin Li
// @receiver r
// @param ctx
// @param arts
// @return error
func (r *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := range arts {

		arts[i].Content = arts[i].CreateAbstract()
	}

	val, err := json.Marshal(&arts)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}
