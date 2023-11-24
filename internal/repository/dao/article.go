// Package dao
// @Description:
package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewGormArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

// @func: Insert
// @date: 2023-11-23 00:59:39
// @brief: 数据库-新建帖子记录
// @author: Kewin Li
// @receiver g
// @param ctx
// @param a
// @return int64
// @return error
func (g *GormArticleDao) Insert(ctx context.Context, art Article) (int64, error) {

	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error

	return art.Id, err
}

type Article struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Title    string `gorm:"type=varchar(256)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
