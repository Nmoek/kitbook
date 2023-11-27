package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDao interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

type GormArticleAuthorDao struct {
	db *gorm.DB
}

func NewGormArticleAuthorDao(db *gorm.DB) *GormArticleAuthorDao {
	return &GormArticleAuthorDao{
		db: db,
	}
}

func (g *GormArticleAuthorDao) Create(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleAuthorDao) Update(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}
