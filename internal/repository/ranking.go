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

	// 以下为双缓存使用 V1
	redisCache *cache.RedisRankingCache
	localCache *cache.LocalRankingCache
}

func NewCacheRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CacheRankingRepository{
		cache: cache,
	}
}

func NewCacheRankingRepositoryV1(redisCache *cache.RedisRankingCache,
	localCache *cache.LocalRankingCache) RankingRepository {
	return &CacheRankingRepository{
		redisCache: redisCache,
		localCache: localCache,
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

// @func: ReplaceTopNV1
// @date: 2023-12-30 22:18:36
// @brief: 热榜服务-热榜数据放入缓存-双缓存设计
// @author: Kewin Li
// @receiver c
// @param ctx
// @param arts
// @return error
func (c *CacheRankingRepository) ReplaceTopNV1(ctx context.Context, arts []domain.Article) error {
	err := c.localCache.Set(ctx, arts)
	if err != nil {
		//TODO: 日志埋点, 本地缓存更新失败
		// 没有必要报错退出
	}

	return c.redisCache.Set(ctx, arts)
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

// @func: GetTopN
// @date: 2023-12-30 22:12:08
// @brief: 热榜服务-热榜数据取出缓存-双缓存设计
// @author: Kewin Li
// @receiver c
// @param ctx
// @return []domain.Article
// @return error
func (c *CacheRankingRepository) GetTopNV1(ctx context.Context) ([]domain.Article, error) {
	// 1. 先查本地缓存
	arts, err := c.localCache.Get(ctx)
	// 注意：这里查询失败有两层含义：本地缓存出错、本地缓存数据过期
	if err == nil {
		return arts, nil
	}

	// 2. 再查redis缓存
	arts, err = c.redisCache.Get(ctx)
	// 2.1 redis查询失败可以考虑直接使用过期的本地缓存数据
	if err != nil {
		return c.localCache.ForceGet(ctx)
	}

	// 3. 本地缓存回写
	err = c.redisCache.Set(ctx, arts)
	if err != nil {
		//TODO: 日志埋点, 本地缓存回写失败
		// 没有必要直接报错返回
	}

	return arts, nil

}
