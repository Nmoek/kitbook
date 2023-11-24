// Package service
// @Description:
package service

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
)

type ArticleServer interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

// NormalArticleService
// @Description: 普通帖子服务(免费、免登录)
type NormalArticleService struct {
	repo repository.ArticleRepository
}

func NewNormalArticleService(repo repository.ArticleRepository) ArticleServer {
	return &NormalArticleService{
		repo: repo,
	}
}

// @func: Save
// @date: 2023-11-23 00:03:50
// @brief: 帖子保存
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
func (n *NormalArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	return n.repo.Create(ctx, art)
}
