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

var ErrUserMismatch = errors.New("帖子ID和用户ID不匹配")

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, artId int64, authorId int64, status uint8) error
	GetByAuthor(ctx context.Context, userId int64, offset int, limit int) ([]Article, error)
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
// @brief: 新建帖子记录
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
			"status":  art.Status,
			"utime":   time.Now().UnixMilli(),
		})

	if result.RowsAffected <= 0 {
		return ErrUserMismatch
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

	// Gorm闭包开启事务
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var err error
		dao := NewGormArticleDao(tx)

		if art.Id > 0 {
			// id存在视为更新
			err = dao.UpdateById(ctx, art)

		} else {
			// id不存在视为新建
			art.Id, err = dao.Insert(ctx, art)

		}
		if err != nil {
			return err
		}

		now := time.Now().UnixMilli()
		art.Utime = now
		art.Ctime = now
		publishArt := PublishedArticle(art)

		// 操作线上表
		// id冲突时处理
		// TODO: 该函数的作用是?
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
				"status":  publishArt.Status,
				"utime":   publishArt.Utime,
			}),
		}).Create(&publishArt).Error

		return err
	})

	return art.Id, err
}

// @func: SyncStatus
// @date: 2023-11-29 00:24:10
// @brief: 帖子状态同步-闭包实现事务
// @author: Kewin Li
// @receiver g
// @param ctx
// @param artId
// @param authorId
// @param status
// @return int64
// @return error
func (g *GormArticleDao) SyncStatus(ctx context.Context, artId int64, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 修改制作库
		res := tx.Model(&Article{}).
			Where("id = ?", artId).
			Where("author_id = ?", authorId).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})

		if res.Error != nil {
			return res.Error
		}

		// 更新无效，说明帖子ID和作者ID不匹配
		if res.RowsAffected <= 0 {
			return ErrUserMismatch
		}

		// 2. 修改线上库
		return tx.Model(&PublishedArticle{}).
			Where("id = ?", artId).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error
	})

}

type Article struct {
	Id       int64  `gorm:"primaryKey, autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(256)" bson:"title,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

// @func: GetByAuthor
// @date: 2023-12-04 00:26:47
// @brief: 帖子服务-查询创作者创作列表
// @author: Kewin Li
// @receiver g
// @param ctx
// @param userId
// @param offset
// @param limit
// @return []domain.Article
// @return error
func (g *GormArticleDao) GetByAuthor(ctx context.Context, userId int64, offset int, limit int) ([]Article, error) {

	var arts []Article
	err := g.db.WithContext(ctx).Where("author_id = ?", userId).
		Offset(offset).
		Limit(limit).Order("utime DESC"). // 最新修改的排在前面
		Find(&arts).Error

	return arts, err

}

// 同步数据-同库不同表 使用衍生类型拓展一张一样的表结构
type PublishedArticle Article
