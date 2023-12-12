package repository

import (
	"context"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
}

type ArticleInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
}

func NewArticleInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache) InteractiveRepository {
	return &ArticleInteractiveRepository{
		dao:   dao,
		cache: cache,
	}
}

// @func: IncreaseReadCnt
// @date: 2023-12-11 23:38:34
// @brief: 阅读数+1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @return error
func (a *ArticleInteractiveRepository) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := a.dao.IncreaseReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}

	// 更新缓存
	// 部分失败问题
	return a.cache.IncreaseReadCntIfPresent(ctx, biz, bizId)
}
