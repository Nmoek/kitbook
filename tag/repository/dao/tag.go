package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type TagDao interface {
	CreateTag(ctx context.Context, tag Tag) (int64, error)
	GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error)
	GetTags(ctx context.Context, offset int, limit int) ([]Tag, error)

	CreateTagBiz(ctx context.Context, tagsBiz []TagBiz) error
	GetTagsByBiz(ctx context.Context, bizId int64, biz string, uid int64) ([]Tag, error)
}

type GormTagDao struct {
	db *gorm.DB
}

func NewGormTagDao(db *gorm.DB) TagDao {
	return &GormTagDao{
		db: db,
	}
}

func (g *GormTagDao) CreateTag(ctx context.Context, tag Tag) (int64, error) {
	now := time.Now().UnixMilli()
	tag.Ctime = now
	tag.Utime = now
	err := g.db.WithContext(ctx).Create(&tag).Error

	return tag.Id, err
}

func (g *GormTagDao) GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error) {
	var res []Tag
	err := g.db.WithContext(ctx).Where("uid = ?", uid).Find(&res).Error
	return res, err
}

func (g *GormTagDao) GetTags(ctx context.Context, offset int, limit int) ([]Tag, error) {
	var res []Tag
	err := g.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (g *GormTagDao) CreateTagBiz(ctx context.Context, tagsBiz []TagBiz) error {
	if len(tagsBiz) <= 0 {
		return nil
	}

	now := time.Now().UnixMilli()
	for _, tag := range tagsBiz {
		tag.Ctime = now
		tag.Utime = now
	}

	return g.db.Transaction(func(tx *gorm.DB) error {
		first := tagsBiz[0]

		// TODO: 这里似乎有并发问题
		// 1. 正常用户并不需要考虑
		// 2. 攻击者能够恶意攻击
		err := tx.Model(&TagBiz{}).Delete("uid = ? AND biz = ? AND biz_id = ?", first.Uid, first.Biz, first.BizId).Error
		if err != nil {
			return err
		}

		return tx.Create(&tagsBiz).Error
	})
}

// @func: GetTagsByBiz
// @date: 2024-03-12 22:25:02
// @brief: 查询某个资源上所有标签信息-两次查询版
// @author: Kewin Li
// @receiver g
// @param ctx
// @param bizId
// @param biz
// @param uid
// @return []Tag
// @return error
func (g *GormTagDao) GetTagsByBiz(ctx context.Context, bizId int64, biz string, uid int64) ([]Tag, error) {
	// 第一次查出 用户给某个资源打上的所有标签ID, 即：Tid
	var bizTags []TagBiz
	err := g.db.WithContext(ctx).Where("uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&bizTags).Error
	if err != nil {
		return nil, err
	}

	// 第二次查出 Tid对应的具体标签信息
	tids := make([]int64, 0, len(bizTags))
	for _, b := range bizTags {
		tids = append(tids, b.Tid)
	}

	var tags []Tag
	err = g.db.WithContext(ctx).Where("tid IN ?", tids).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// @func: GetTagsByBizV2
// @date: 2024-03-12 22:24:56
// @brief: 查询某个资源上所有标签信息-Join查询版
// @author: Kewin Li
// @receiver g
// @param ctx
// @param bizId
// @param biz
// @param uid
// @return []Tag
// @return error
func (g *GormTagDao) GetTagsByBizV2(ctx context.Context, bizId int64, biz string, uid int64) ([]Tag, error) {

	var bizTags []TagBiz
	// 注意: 同名字段需要指明是哪个结构体
	err := g.db.WithContext(ctx).Model(&TagBiz{}).InnerJoins("Tag", g.db.Model(&Tag{})).
		Where("Tag.uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&bizTags).Error
	if err != nil {
		return nil, err
	}

	res := make([]Tag, 0, len(bizTags))
	for _, b := range bizTags {
		res = append(res, *b.Tag)
	}

	return res, nil
}

// Tag
// @Description: 标签表结构定义
type Tag struct {
	Id   int64  `gorm:"primaryKey, autoIncrement"`
	Uid  int64  `gorm:"uniqueIndex:uid_name_idx"`
	Name string `gorm:"type:varchar(4096);uniqueIndex:uid_name_idx"`

	Ctime int64
	Utime int64
}

// TagBiz
// @Description: 资源标签表结构定义
type TagBiz struct {
	Id    int64  `gorm:"primaryKey, autoIncrement"`
	BizId int64  `gorm:"index:biz_type_id"`
	Biz   string `gorm:"index:biz_type_id"`

	Uid int64 `gorm:"index"`
	Tid int64
	// 联合删除
	Tag   *Tag `gorm:"ForeignKey:Tid;AssociationForeignKey:Id;constraint:Ondelete:CASCADE"`
	Ctime int64
	Utime int64
}
