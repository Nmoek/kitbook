package service

import (
	"context"
	"kitbook/internal/repository"
)

type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
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
