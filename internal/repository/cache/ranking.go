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
	Get(ctx context.Context) ([]domain.Article, error)
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

// @func: Get
// @date: 2023-12-30 21:48:23
// @brief:热榜缓存-取出数据
// @author: Kewin Li
// @receiver r
// @param ctx
// @return []domain.Article
// @return error
func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	var arts []domain.Article

	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &arts)
	return arts, err
}
