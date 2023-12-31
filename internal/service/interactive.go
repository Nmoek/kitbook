package service

import (
	"context"
	"golang.org/x/sync/errgroup"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/pkg/logger"
)

type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, userId int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error
	Collect(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	CancelCollect(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error)
	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type ArticleInteractiveService struct {
	repo repository.InteractiveRepository

	l logger.Logger
}

func NewArticleInteractiveService(repo repository.InteractiveRepository, l logger.Logger) InteractiveService {
	return &ArticleInteractiveService{
		repo: repo,
		l:    l,
	}
}

// @func: IncreaseReadCnt
// @date: 2023-12-11 23:39:23
// @brief: 阅读数+1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizIf
// @return error
func (a *ArticleInteractiveService) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	return a.repo.IncreaseReadCnt(ctx, biz, bizId)
}

// @func: Like
// @date: 2023-12-13 21:57:28
// @brief: 点赞数+1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param artId
// @param userId
// @return error
func (a *ArticleInteractiveService) Like(ctx context.Context, biz string, bizId int64, userId int64) error {
	return a.repo.IncreaseLikeCnt(ctx, biz, bizId, userId)
}

// @func: CancelLike
// @date: 2023-12-13 21:57:37
// @brief: 点赞数-1
// @author: Kewin Li
// @receiver a
// @param ctx
// @param artId
// @param userId
// @return error
func (a *ArticleInteractiveService) CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error {
	return a.repo.DecreaseLikeCnt(ctx, biz, bizId, userId)
}

// @func:
// @date: 2023-12-14 01:45:00
// @brief: 收藏帖子
// @author: Kewin Li
// @receiver a
// @return unc
func (a *ArticleInteractiveService) Collect(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	return a.repo.IncreaseCollectItem(ctx, biz, bizId, collectId, userId)
}

// @func:
// @date: 2023-12-14 01:45:00
// @brief: 取消收藏帖子
// @author: Kewin Li
// @receiver a
// @return unc
func (a *ArticleInteractiveService) CancelCollect(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	return a.repo.DecreaseCollectItem(ctx, biz, bizId, collectId, userId)
}

// @func: Get
// @date: 2023-12-14 13:04:12
// @brief: 聚合查询-阅读数、点赞数、收藏数
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return domain.Interactive
// @return error
func (a *ArticleInteractiveService) Get(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error) {
	// 1. 查询 互动内容数据
	intr, err := a.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}

	var eg errgroup.Group

	// 2. 并发查询 是否点过赞
	eg.Go(func() error {
		var err2 error
		intr.Liked, err2 = a.repo.Liked(ctx, biz, bizId, userId)
		if err2 != nil {
			a.l.ERROR("查询用户点赞信息失败",
				logger.Error(err),
				logger.Field{"biz", biz},
				logger.Int[int64]("biz_id", bizId),
				logger.Int[int64]("user_id", userId))
			return err2
		}
		return nil
	})

	// 3. 并发查询 是否收藏过
	eg.Go(func() error {
		var err2 error
		intr.Collected, err2 = a.repo.Collectd(ctx, biz, bizId, userId)
		if err2 != nil {
			a.l.ERROR("查询用户收藏信息失败",
				logger.Error(err),
				logger.Field{"biz", biz},
				logger.Int[int64]("biz_id", bizId),
				logger.Int[int64]("user_id", userId))
			return err2
		}

		return nil
	})
	// TODO: 弱校验, 互动数据查询失败对于文章主体内容并不影响，不一定非要报错处理
	// TODO: 系统降级，当负载较高时, 上述两个并发查询都可以不再进行查询
	return intr, eg.Wait()
}

// @func: GetByIds
// @date: 2023-12-27 00:49:52
// @brief: 热榜服务-分批查出点赞数
// @author: Kewin Li
// @receiver a
// @param ctx
// @param biz
// @param bizIds
// @return map[int64]domain.Article
// @return error
func (a *ArticleInteractiveService) GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	res := map[int64]domain.Interactive{}
	intrs, err := a.repo.GetByIds(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}

	for _, intr := range intrs {
		res[intr.BizId] = intr
	}

	return res, nil
}
