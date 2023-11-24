// Package dao
// @Description:
package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrInvalidUpdate = errors.New("帖子ID或用户ID不匹配")

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
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

// @func: UpdateById
// @date: 2023-11-24 21:02:25
// @brief:  数据库-修改帖子记录按Id
// @author: Kewin Li
// @receiver g
// @param ctx
// @param article
// @return error
func (g *GormArticleDao) UpdateById(ctx context.Context, art Article) error {

	result := g.db.WithContext(ctx).Model(&Article{}).
		Where("id = ?", art.Id).
		Where("author_id = ?", art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   time.Now().UnixMilli(),
		})

	if result.RowsAffected <= 0 {
		return ErrInvalidUpdate
	}

	return result.Error

}

type Article struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Title    string `gorm:"type=varchar(256)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
