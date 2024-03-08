package repository

import (
	"context"
	"kitbook/search/domain"
	"kitbook/search/repository/dao"
)

type articleSyncRepository struct {
	dao dao.ArticleDao
}

func NewArticleSyncRepository(dao dao.ArticleDao) ArticleSyncRepository {
	return &articleSyncRepository{
		dao: dao,
	}
}

func (a *articleSyncRepository) InputArticle(ctx context.Context, art domain.Article) error {
	return a.dao.InputArticle(ctx, a.convertsArtDao(&art))
}

func (a *articleSyncRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (u *articleSyncRepository) convertsArtDao(art *domain.Article) dao.Article {
	return dao.Article{}
}
