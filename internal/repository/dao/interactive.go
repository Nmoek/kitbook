package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
}

type GORMInteractiveDao struct {
	db *gorm.DB
}

func NewGORMInteractiveDao(db *gorm.DB) InteractiveDao {
	return &GORMInteractiveDao{
		db: db,
	}
}

// @func: IncreaseReadCnt
// @date: 2023-12-11 00:15:53
// @brief: 数据库操作-阅读数+1(UpSert语义)
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param id
// @return error
func (g *GORMInteractiveDao) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()

	return g.db.WithContext(ctx).Clauses(&clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			// Expr可将SQL用于表达式传递
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

// Interactive
// @Description: 阅读数、点赞数、收藏数三合一
type Interactive struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	// TODO: 建立联合唯一索引<bizId, biz>
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"type=varchar(128), uniqueIndex:biz_type_id"`

	// 阅读数
	ReadCnt int64
	// 点赞数
	LikeCnt int64
	// 收藏数
	CollectCnt int64
	Utime      int64
	Ctime      int64
}
