package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type CommentDao interface {
	Insert(ctx context.Context, cmt Comment) error
	Delete(ctx context.Context, cmt Comment) error
}

type GormCommentDao struct {
	db *gorm.DB
}

func NewGormCommentDao(db *gorm.DB) CommentDao {
	return &GormCommentDao{db: db}
}

func (g *GormCommentDao) Insert(ctx context.Context, cmt Comment) error {
	return g.db.WithContext(ctx).Create(&cmt).Error
}

func (g *GormCommentDao) Delete(ctx context.Context, cmt Comment) error {
	return g.db.WithContext(ctx).Delete(&Comment{
		Id: cmt.Id,
	}).Error
}

// Comment
// @Description: 评论表
type Comment struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 大数据会形成用户画像
	Uid int64

	BizId   int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz     string `gorm:"uniqueIndex:biz_type_id"`
	Content string `gorm:"type:text"`
	//根节点评论, 解决加载更多评论问题
	RootId sql.NullInt64 `gorm:"index"`
	//父节点评论
	Pid   sql.NullInt64 `gorm:"index;Foreignkey:Pid;AssociationForeignkey:Id;constraint:Ondelete:CASCADE"`
	Utime int64
	Ctime int64
}
