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
	GetLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) (UserLikeInfo, error)

	AddCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	DelCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error
	GetCollectionItem(ctx context.Context, biz string, bizId int64, userId int64) (UserCollectInfo, error)
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	BatchIncreaseReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
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

// @func: BatchIncreaseReadCnt
// @date: 2023-12-19 03:20:46
// @brief: 数据库操作-一次事务提交多条记录
// @author: Kewin Li
// @receiver g
// @param ctx
// @param bizs
// @param bizIds
// @return error
func (g *GORMInteractiveDao) BatchIncreaseReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {

	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDao := NewGORMInteractiveDao(tx)

		for i := 0; i < len(bizs); i++ {
			err := txDao.IncreaseReadCnt(ctx, bizs[i], bizIds[i])
			if err != nil {
				//TODO: 日志埋点
				return err
			}
		}

		return nil
	})
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
			Where("user_id = ? AND biz_id = ? AND biz = ?", userId, bizId, biz).
			Updates(map[string]any{
				"status": 0, //当前点赞无效
				"utime":  now,
			}).Error

		if err != nil {
			// TODO: 日志埋点
		}

		// 2. 互动表，点赞数-1
		return tx.Model(&Interactive{}).
			Where("biz_id = ? AND biz = ?", bizId, biz).
			Updates(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"utime":    now,
			}).Error
	})
}

// @func: GetLikeInfo
// @date: 2023-12-15 16:52:28
// @brief: 查询用户点赞信息
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return UserLikeInfo
// @return error
func (g *GORMInteractiveDao) GetLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) (UserLikeInfo, error) {
	var info UserLikeInfo
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ? AND user_id = ?", biz, bizId, userId).
		First(&info).Error
	return info, err
}

// @func: AddCollectionItem
// @date: 2023-12-14 02:03:10
// @brief: 增加收藏内容
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param collectId
// @param userId
// @return error
func (g *GORMInteractiveDao) AddCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	now := time.Now().UnixMilli()

	// 1. 收藏信息表
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status":     1,
				"collect_id": collectId,
				"utime":      now,
			}),
		}).Create(&UserCollectInfo{
			UserId:    userId,
			BizId:     bizId,
			Biz:       biz,
			Status:    1,
			CollectId: collectId,
			Utime:     now,
			Ctime:     now,
		}).Error

		if err != nil {
			//TODO: 日志埋点
			return err
		}

		// 2. 互动表
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			BizId:      bizId,
			Biz:        biz,
			CollectCnt: 1,
			Utime:      now,
			Ctime:      now,
		}).Error
	})
}

// @func: DelCollectionItem
// @date: 2023-12-14 02:03:20
// @brief: 去除收藏内容
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param collectId
// @param userId
// @return error
func (g *GORMInteractiveDao) DelCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	now := time.Now().UnixMilli()

	// 1. 软删除 收藏信息
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserCollectInfo{}).
			Where("user_id = ? AND biz_id = ? AND biz = ?", userId, bizId, biz).
			Updates(map[string]any{
				"status": 0,
				"utime":  now,
			}).Error

		if err != nil {
			return err
		}
		// 2. 互动表 收藏数-1
		return tx.Model(&Interactive{}).
			Where("user_id = ? AND biz_id = ? AND biz = ? AND status = ? ", userId, bizId, biz, 1). //注意: status=1 才是有效
			Updates(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt` - 1"),
				"utime":       now,
			}).Error
	})

}

// @func: GetCollectionItem
// @date: 2023-12-15 16:58:12
// @brief: 查询用户收藏信息
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return UserCollectInfo
func (g *GORMInteractiveDao) GetCollectionItem(ctx context.Context, biz string, bizId int64, userId int64) (UserCollectInfo, error) {
	var item UserCollectInfo
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ? AND user_id = ? AND status = ?", bizId, biz, userId, 1). //注意: status=1 才是有效
		First(&item).Error
	return item, err
}

// @func: Get
// @date: 2023-12-15 17:37:20
// @brief: 查询用户互动信息
// @author: Kewin Li
// @receiver g
// @param ctx
// @param biz
// @param bizId
// @return Interactive
// @return error
func (g *GORMInteractiveDao) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var intr Interactive
	err := g.db.WithContext(ctx).Where("biz_id = ? AND biz = ?", bizId, biz).First(&intr).Error
	return intr, err
}

// Interactive
// @Description: 阅读数、点赞数、收藏数三合一
type Interactive struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	// 建立联合唯一索引<bizId, biz>
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`

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
// @Description: 用户已点赞的信息表, 记录当前用户给什么点了赞
type UserLikeInfo struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 以用户ID为主字段查询
	UserId int64 `gorm:"uniqueIndex:uid_biz_type_id"`
	// BizId + Biz 共同表示哪个业务的哪一条记录
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`

	// 点赞是否有效
	Status int8
	Utime  int64
	Ctime  int64
}

// UserCollectInfo
// @Description: 用户已收藏的信息表, 记录当前哪个帖子被收藏在哪个收藏夹中

type UserCollectInfo struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	// 以用户ID为主字段查询
	UserId int64 `gorm:"uniqueIndex:uid_biz_type_id"`
	// BizId + Biz 共同表示哪个业务的哪一条记录
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`

	// 被收藏在哪一个收藏夹
	// 注意: 一个资源只能被收藏一次, 否则就会出现多个收藏夹中同一个资源
	CollectId int64 `gorm:"index"`
	// 收藏是否还有效
	Status int8
	Utime  int64
	Ctime  int64
}
