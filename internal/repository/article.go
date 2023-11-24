package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDao
}

func NewCacheArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, ConvertsDaoArticle(&art))
}

// @func: convertsDominUser
// @date: 2023-10-09 02:08:11
// @brief: 转化为domin的Article结构体
// @author: Kewin Li
// @param user
// @return domain.User
func ConvertsDomainArticle(art *dao.Article) domain.Article {
	return domain.Article{}
}

// @func: ConvertsDaoUser
// @date: 2023-11-23 00:55:51
// @brief: 转化为dao的Article结构体
// @author: Kewin Li
// @param user
// @return dao.User
func ConvertsDaoArticle(art *domain.Article) dao.Article {
	return dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
