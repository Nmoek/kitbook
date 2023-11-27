// Package dao
// @Description:
package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrInvalidUpdate = errors.New("帖子ID或用户ID不匹配")

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
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

// @func: Sync
// @date: 2023-11-27 12:48:03
// @brief: 帖子发表-数据同步-用户控制事务
// @author: Kewin Li
// @receiver g
// @param ctx
// @param article
// @return int64
// @return error
func (g *GormArticleDao) SyncV1(ctx context.Context, art Article) (int64, error) {
	var id int64
	var err error

	// 开启事务
	dbCtx := g.db.WithContext(ctx).Begin()
	if dbCtx.Error != nil {
		return 0, dbCtx.Error
	}
	defer dbCtx.Rollback() //事务回滚

	authorDao := NewGormArticleAuthorDao(dbCtx)
	readerDao := NewGormArticleReaderDao(dbCtx)

	if art.Id > 0 {
		err = authorDao.Update(ctx, art)
		if err != nil {
			return art.Id, err
		}

	} else {

		id, err = authorDao.Create(ctx, art)
		if err != nil {
			return 0, err
		}
		art.Id = id
	}

	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	publishArt := PublishedArticle(art)

	// 操作另外一张表
	// id冲突时处理
	err = readerDao.db.Clauses(clause.OnConflict{
		// MYSQL:
		// INSERT xxx DUPLICATE KEY SET `title` = ?

		//其他方言:
		//sqlite: INSERT xxx ON CONFLICT DO UPDATES WHERE xxx

		// Columns兼容其他方言
		Columns: []clause.Column{{Name: "id"}},
		// TODO: Mysql仅支持该字段, 其他不生效
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   publishArt.Title,
			"content": publishArt.Content,
			"utime":   publishArt.Utime,
		}),
	}).Create(&publishArt).Error

	err = dbCtx.Commit().Error
	return art.Id, err
}

// @func: Sync
// @date: 2023-11-27 18:06:37
// @brief: 帖子发表-数据同步-闭包控制事务
// @author: Kewin Li
// @receiver g
// @param ctx
// @param art
// @return int64
// @return error
func (g *GormArticleDao) Sync(ctx context.Context, art Article) (int64, error) {
	id := art.Id
	// 开启事务
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var err error
		dao := NewGormArticleAuthorDao(tx)

		if id > 0 {
			err = dao.Update(ctx, art)

		} else {

			id, err = dao.Create(ctx, art)

			art.Id = id
		}
		if err != nil {
			return err
		}

		now := time.Now().UnixMilli()
		art.Utime = now
		art.Ctime = now
		publishArt := PublishedArticle(art)

		// 操作另外一张表
		// id冲突时处理
		err = tx.Clauses(clause.OnConflict{
			// MYSQL:
			// INSERT xxx DUPLICATE KEY SET `title` = ?

			//其他方言:
			//sqlite: INSERT xxx ON CONFLICT DO UPDATES WHERE xxx

			// Columns兼容其他方言
			Columns: []clause.Column{{Name: "id"}},
			// TODO: Mysql仅支持该字段, 其他不生效
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   publishArt.Title,
				"content": publishArt.Content,
				"utime":   publishArt.Utime,
			}),
		}).Create(&publishArt).Error

		return err
	})

	return id, err
}

type Article struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Title    string `gorm:"type=varchar(256)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}

// 同步数据-同库不同表 使用衍生类型拓展一张一样的表结构
type PublishedArticle Article
