package repository

import (
	"context"
	"kitbook/search/domain"
	"kitbook/search/repository/dao"
)

type articleSyncRepository struct {
	dao dao.ArticleDao
	tag dao.TagDao
}

func NewArticleSyncRepository(dao dao.ArticleDao, tag dao.TagDao) ArticleRepository {
	return &articleSyncRepository{
		dao: dao,
		tag: tag,
	}
}

func (a *articleSyncRepository) InputArticle(ctx context.Context, art domain.Article) error {
	return a.dao.InputArticle(ctx, a.convertsArtDao(&art))
}

func (a *articleSyncRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {
	artIds, err := a.tag.SearchTag(ctx, uid, "article", keywords)
	if err != nil {
		return nil, err
	}

	res, err := a.dao.SearchArticle(ctx, artIds, keywords)
	return a.convertsArtsDomain(res), err
}

func (a *articleSyncRepository) convertsArtDao(art *domain.Article) dao.Article {
	return dao.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status,
	}
}

func (a *articleSyncRepository) convertsArtsDao(arts []domain.Article) []dao.Article {
	res := make([]dao.Article, 0, len(arts))
	for _, art := range arts {
		res = append(res, dao.Article{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Status:  art.Status,
		})
	}
	return res
}

func (a *articleSyncRepository) convertsArtDomain(art *dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  art.Status,
	}
}

func (a *articleSyncRepository) convertsArtsDomain(arts []dao.Article) []domain.Article {
	res := make([]domain.Article, 0, len(arts))
	for _, art := range arts {
		res = append(res, domain.Article{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Status:  art.Status,
		})
	}
	return res
}
