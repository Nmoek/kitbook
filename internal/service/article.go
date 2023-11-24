// Package service
// @Description:
package service

import (
	"context"
	"errors"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
)

var ErrInvalidUpdate = errors.New("非法操作")

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
	if art.Id > 0 {
		err := n.repo.Update(ctx, art)
		if err == repository.ErrInvalidUpdate {
			return -1, ErrInvalidUpdate
		}
		return art.Id, err
	}
	return n.repo.Create(ctx, art)
}
