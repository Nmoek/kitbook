package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type RewardDao interface {
	Insert(ctx context.Context, rwd Reward) (int64, error)
	Update(ctx context.Context, rwd Reward) error
	UpdateStatus(ctx context.Context, rid int64, status uint8) error
	GetByRid(ctx context.Context, rid int64) (Reward, error)
}

type GormRewardDao struct {
	db *gorm.DB
}

func NewGormRewardDao(db *gorm.DB) RewardDao {
	return &GormRewardDao{
		db: db,
	}
}

func (g *GormRewardDao) Insert(ctx context.Context, rwd Reward) (int64, error) {
	err := g.db.WithContext(ctx).Create(&rwd).Error
	return rwd.Id, err
}

func (g *GormRewardDao) Update(ctx context.Context, rwd Reward) error {
	return g.db.WithContext(ctx).Model(&Reward{}).
		Where("biz_id = ? AND biz = ?", rwd.BizId, rwd.Biz).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": rwd.Status,
		}).Error
}

func (g *GormRewardDao) UpdateStatus(ctx context.Context, rid int64, status uint8) error {

	return g.db.WithContext(ctx).Model(&Reward{}).
		Where("id = ? ", rid).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": status,
		}).Error
}

func (g *GormRewardDao) GetByRid(ctx context.Context, rid int64) (Reward, error) {
	var res Reward
	err := g.db.WithContext(ctx).Where("id = ?", rid).First(&res).Error
	return res, err
}

// Reward
// @Description: 打赏信息表
type Reward struct {
	Id      int64  `gorm:"primaryKey, autoIncrement"`
	BizId   int64  `gorm:"uniqueIndex:rwd_biz_type_id"`
	Biz     string `gorm:"uniqueIndex:rwd_biz_type_id"`
	BizName string `gorm:"type:varchar(2048)"`
	// 打赏人和被打赏人不可能是唯一的
	TargetUid int64 `gorm:"index:rwd_uid_type_id"` // 被打赏人ID
	Uid       int64 `gorm:"index:rwd_uid_type_id"` // 打赏人ID
	Amt       int64 // 打赏金额
	Status    uint8

	Utime int64
	Ctime int64
}
