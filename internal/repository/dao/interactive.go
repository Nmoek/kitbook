package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDao interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	AddLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	DelLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
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

// @func: AddLikeInfo
// @date: 2023-12-13 22:10:12
// @brief: 添加用户点赞信息
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (g *GORMInteractiveDao) AddLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()

	// 1. 用户点赞表 增加信息
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeInfo{
			UserId: userId,
			BizId:  bizId,
			Biz:    biz,
			Status: 1,
			Utime:  now,
			Ctime:  now,
		}).Error

		if err != nil {
			// TODO: 日志埋点
		}

		// 2. 互动表，点赞数+1
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			BizId:   bizId,
			Biz:     biz,
			LikeCnt: 1,
			Utime:   now,
			Ctime:   now,
		}).Error

	})
}

// @func: DelLikeInfo
// @date: 2023-12-13 22:10:23
// @brief: 删除用户点赞信息
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (g *GORMInteractiveDao) DelLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()

	// 1. 用户点赞表 软删除
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeInfo{}).
			Where("user_id = ? AND biz_id = ? AND biz = ?").
			Updates(map[string]any{
				"status": 0,
				"utime":  now,
			}).Error

		if err != nil {
			// TODO: 日志埋点
		}

		// 2. 互动表，点赞数-1
		return tx.Model(&Interactive{}).
			Where("biz_id = ? AND biz = ?").
			Updates(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"utime":    now,
			}).Error
	})
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

// UserLikeInfo
// @Description: 用户被点赞的信息表
type UserLikeInfo struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 以用户ID为主字段查询
	UserId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"uniqueIndex:uid_biz_type_id"`

	// 点赞是否有效
	Status int8
	Utime  int64
	Ctime  int64
}
