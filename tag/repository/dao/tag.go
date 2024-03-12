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

// Tag
// @Description: 标签表结构定义
type Tag struct {
	Id   int64  `gorm:"primaryKey, autoIncrement"`
	Uid  int64  `gorm:"uniqueIndex:uid_name_idx"`
	Name string `gorm:"type:varchar(4096);uniqueIndex:uid_name_idx"`

	Ctime int64
	Utime int64
}
