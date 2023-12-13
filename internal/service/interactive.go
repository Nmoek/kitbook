package service

import (
	"context"
	"kitbook/internal/repository"
)

type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, userId int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error
}

type ArticleInteractiveService struct {
	repo repository.InteractiveRepository
}

func NewArticleInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &ArticleInteractiveService{
		repo: repo,
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
