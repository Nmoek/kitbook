package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

type CacheRankingRepository struct {
	cache cache.RankingCache
}

func NewCacheRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CacheRankingRepository{
		cache: cache,
	}
}

// @func: ReplaceTopN
// @date: 2023-12-29 22:58:53
// @brief: 热榜服务-热榜数据放入缓存
// @author: Kewin Li
// @receiver c
// @param ctx
func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return c.cache.Set(ctx, arts)
}
