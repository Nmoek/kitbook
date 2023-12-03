// Package dao
// @Description: S3(OSS)实现数据存储
package dao

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

type S3ArticleDAO struct {
	GormArticleDao
	oss *s3.S3
}

func NewS3ArticleDAO(db *gorm.DB, oss *s3.S3) *S3ArticleDAO {
	return &S3ArticleDAO{
		GormArticleDao: GormArticleDao{db},
		oss:            oss,
	}
}

// @func: Sync
// @date: 2023-12-03 22:29:54
// @brief: 帖子发表-数据同步-S3实现
// @author: Kewin Li
// @receiver s
// @param ctx
// @param art
// @return int64
// @return error
func (s *S3ArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// Gorm闭包开启事务
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

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
		publishArt := PublishedArticleV2{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Ctime:    now,
			Utime:    now,
			Status:   art.Status,
		}

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
				"title":  publishArt.Title,
				"status": publishArt.Status,
				"utime":  publishArt.Utime,
			}),
		}).Create(&publishArt).Error

		return err
	})

	if err != nil {
		return art.Id, err
	}

	_, err = s.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("kitbook-1314583317"),
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})

	return art.Id, err
}

// @func: SyncStatus
// @date: 2023-12-03 22:42:30
// @brief: 帖子状态同步-S3实现
// @author: Kewin Li
// @receiver s
// @param ctx
// @param artId
// @param authorId
// @param status
// @return error
// TODO: 部署网站进行验证
func (s *S3ArticleDAO) SyncStatus(ctx *gin.Context, artId int64, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		return tx.Model(&PublishedArticleV2{}).
			Where("id = ?", artId).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error
	})

	if err != nil {
		return err
	}

	const ArticleStatusPrivate = 3
	if status == ArticleStatusPrivate {
		_, err = s.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: ekit.ToPtr[string]("kitbook-1314583317"),
			Key:    ekit.ToPtr[string](strconv.FormatInt(artId, 10)),
		})

	}

	return err
}

// PublishedArticleV2
// @Description: OSS方案
type PublishedArticleV2 struct {
	Id    int64  `gorm:"primaryKey, autoIncrement" bson:"id,omitempty"`
	Title string `gorm:"type=varchar(256)" bson:"title,omitempty"`
	// 内容不存储
	//Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8 `bson:"status,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	Utime    int64 `bson:"utime,omitempty"`
}
