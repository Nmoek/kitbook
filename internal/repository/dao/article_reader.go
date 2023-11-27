package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDao interface {
	Upsert(crx context.Context, art Article) error
	UpsertV2(crx context.Context, art PublishedArticle) error
}

type GormArticleReaderDao struct {
	db *gorm.DB
}

func NewGormArticleReaderDao(db *gorm.DB) *GormArticleReaderDao {
	return &GormArticleReaderDao{
		db: db,
	}
}

func (g *GormArticleReaderDao) Upsert(crx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormArticleReaderDao) UpsertV2(crx context.Context, art PublishedArticle) error {
	//TODO implement me
	panic("implement me")
}
