package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GormFollowDao struct {
	db *gorm.DB
}

func NewGormFollowDao(db *gorm.DB) FollowDao {
	return &GormFollowDao{db: db}
}

func (g *GormFollowDao) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	now := time.Now().UnixMilli()
	f.Utime = now
	f.Ctime = now
	f.Status = FollowRelationStatusActive
	// 保持UPSERT语义
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"status": FollowRelationStatusActive,
			"utime":  now,
		}),
	}).Create(&f).Error

}

func (g *GormFollowDao) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	return g.db.WithContext(ctx).Where("followee = ? AND follower = ? ", followee, follower).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": status,
		}).Error
}

func (g *GormFollowDao) FollowRelationList(ctx context.Context, follower int64, offset int64, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := g.db.WithContext(ctx).Where("follower = ? AND status = ?", follower, FollowRelationStatusActive).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (g *GormFollowDao) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	// 特别注意： 这里followee字段无法命中索引, 属于索引失效范畴
	// 可以考虑给followee单独建索引

	// 查出粉丝中包含uid的所有用户总个数
	err := g.db.WithContext(ctx).
		Select("count(follower)").
		Where("followee = ? AND status = ?", uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GormFollowDao) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	// 查出被关注者名单中包含uid的所有用户总个数
	err := g.db.WithContext(ctx).
		Select("count(followee)").
		Where("follower = ? AND status = ?", uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (g *GormFollowDao) FollowRelationDetail(ctx context.Context, followee int64, follower int64) (FollowRelation, error) {
	var res FollowRelation
	err := g.db.WithContext(ctx).Where("followee = ? AND follower = ? AND status = ?", followee, follower, FollowRelationStatusActive).
		First(&res).Error
	return res, err
}
