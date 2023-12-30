package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
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

// @func: GetTopN
// @date: 2023-12-30 21:47:36
// @brief: 热榜服务-热榜数据取出缓存
// @author: Kewin Li
// @receiver c
// @param ctx
// @return []domain.Article
// @return error
func (c *CacheRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return c.cache.Get(ctx)
}
