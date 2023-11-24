package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
)

var ErrInvalidUpdate = dao.ErrInvalidUpdate

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type CacheArticleRepository struct {
	dao dao.ArticleDao
}

func NewCacheArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

// @func: Create
// @date: 2023-11-24 21:01:34
// @brief: 新建帖子并保存
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, ConvertsDaoArticle(&art))
}

// @func: Update
// @date: 2023-11-24 21:01:45
// @brief: 修改帖子
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return error
func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, ConvertsDaoArticle(&art))
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
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
