package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type CommentDao interface {
	Insert(ctx context.Context, cmt Comment) error
	Delete(ctx context.Context, cmt Comment) error
	FindByBiz(ctx context.Context, bizId int64, biz string, minId int64, limit int64) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, offset int64, limit int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rootId int64, maxId int64, limit int64) ([]Comment, error)
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

// @func: FindByBiz
// @date: 2024-02-11 14:41:51
// @brief: 查询一级评论
// @author: Kewin Li
// @receiver g
// @param ctx
// @param bizId
// @param biz
// @param minId
// @param limit
// @return []Comment
// @return error
func (g *GormCommentDao) FindByBiz(ctx context.Context, bizId int64, biz string, minId int64, limit int64) ([]Comment, error) {
	var res []Comment
	// 查出根结点
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ? AND id < ? AND pid IS NULL", bizId, biz, minId).
		Find(&res).Error
	return res, err
}

// @func: FindRepliesByPid
// @date: 2024-02-11 14:42:00
// @brief: 查询二级评论
// @author: Kewin Li
// @receiver g
// @param ctx
// @param pid
// @param offset
// @param limit
// @return []Comment
// @return error
func (g *GormCommentDao) FindRepliesByPid(ctx context.Context, pid int64, offset int64, limit int64) ([]Comment, error) {
	var res []Comment
	// 二级评论 先加载评论时间最新的
	err := g.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DECRE").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

// @func: FindRepliesByRid
// @date: 2024-02-11 16:47:23
// @brief: 加载更多二级评论
// @author: Kewin Li
// @receiver g
// @param ctx
// @param rootId
// @param maxId
// @param limit
// @return []Comment
// @return error
func (g *GormCommentDao) FindRepliesByRid(ctx context.Context, rootId int64, maxId int64, limit int64) ([]Comment, error) {
	var res []Comment

	err := g.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rootId, maxId).
		Order("id").
		Limit(int(limit)).
		Find(&res).Error
	return res, err
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
