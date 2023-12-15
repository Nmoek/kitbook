package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/pkg/logger"
)

type InteractiveRepository interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	IncreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	DecreaseCollectItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
	Collectd(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
}

type ArticleInteractiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
	l     logger.Logger
}

func NewArticleInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache, l logger.Logger) InteractiveRepository {
	return &ArticleInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
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

// @func: Get
// @date: 2023-12-14 13:07:48
// @brief: 聚合查询
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @return domain.Interactive
// @return error
func (a *ArticleInteractiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {

	// 1. 查缓存
	intrDomain, err := a.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intrDomain, err
	}
	//TODO: 日志埋点，缓存查询失败错误信息

	// 2. 没缓存，查库 查互动信息
	intrDao, err := a.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}

	// 3. 缓存回写
	intrDomain = a.ConvertsDomainInteractive(&intrDao)
	err = a.cache.Set(ctx, biz, bizId, intrDomain)
	if err != nil {
		//TODO: 日志埋点，不一定返回错误
	}

	return intrDomain, nil

}

// @func: Liked
// @date: 2023-12-15 16:42:26
// @brief: 某个业务的某条记录是否被当前用户点赞过
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return bool
// @return error
func (a *ArticleInteractiveRepository) Liked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	_, err := a.dao.GetLikeInfo(ctx, biz, bizId, userId)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

// @func: Collectd
// @date: 2023-12-15 16:42:53
// @brief: 某个业务的某条记录是否被当前用户收藏过
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return bool
// @return error
func (a *ArticleInteractiveRepository) Collectd(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	_, err := a.dao.GetCollectionItem(ctx, biz, bizId, userId)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

// @func: ConvertsDomainInteractive
// @date: 2023-12-15 17:30:22
// @brief: Interactive DAO--->Domain
// @author: Kewin Li
// @receiver a
// @param i
// @return domain.Interactive
func (a *ArticleInteractiveRepository) ConvertsDomainInteractive(i *dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    i.ReadCnt,
		LikeCnt:    i.LikeCnt,
		CollectCnt: i.CollectCnt,
	}
}

// @func: ConvertsDomainInteractive
// @date: 2023-12-15 17:30:22
// @brief: Interactive Domain--->DAO
// @author: Kewin Li
// @receiver a
// @param i
// @return domain.Interactive
func (a *ArticleInteractiveRepository) ConvertsDaoInteractive(i *domain.Interactive) dao.Interactive {
	return dao.Interactive{
		ReadCnt:    i.ReadCnt,
		LikeCnt:    i.LikeCnt,
		CollectCnt: i.CollectCnt,
	}
}
