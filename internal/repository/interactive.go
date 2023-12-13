package repository

import (
	"context"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	IncreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	DecreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
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

// @func: IncreaseLikeCnt
// @date: 2023-12-13 22:01:27
// @brief: 点赞数+1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (a *ArticleInteractiveRepository) IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 1. 数据库
	err := a.dao.AddLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}

	// 2. 更新缓存
	return a.cache.IncreaseLikeCntIfPresent(ctx, biz, bizId)
}

// @func: DecreaseLikeCnt
// @date: 2023-12-13 22:01:32
// @brief: 点赞数-1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (a *ArticleInteractiveRepository) DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 1. 数据库
	err := a.dao.DelLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}

	return a.cache.DecreaseLikeCntIfPresent(ctx, biz, bizId)
}

// @func: IncreaseCollectCnt
// @date: 2023-12-14 01:47:46
// @brief: 收藏帖子
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (a *ArticleInteractiveRepository) IncreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {

	// 1. 修改收藏表、互动表
	err := a.dao.AddCollectionItem(ctx, biz, bizId, collectId, userId)
	if err != nil {
		return err
	}

	// 2. 更新缓存
	return a.cache.IncrCollectionCntIfPresent(ctx, biz, bizId)
}

// @func: DecreaseCollectCnt
// @date: 2023-12-14 01:47:53
// @brief: 取消收藏帖子
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (a *ArticleInteractiveRepository) DecreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	// 1. 修改收藏表、互动表
	err := a.dao.DelCollectionItem(ctx, biz, bizId, collectId, userId)
	if err != nil {
		return err
	}

	// 2. 更新缓存
	return a.cache.DecrCollectionCntIfPresent(ctx, biz, bizId)
}
