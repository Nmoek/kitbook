package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountDao interface {
	AddActivities(ctx context.Context, activities ...AccountActivity) error
}

type GormAccountDao struct {
	db *gorm.DB
}

func NewGormAccountDao(db *gorm.DB) AccountDao {
	return &GormAccountDao{
		db: db,
	}
}

func (g *GormAccountDao) AddActivities(ctx context.Context, activities ...AccountActivity) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		for _, act := range activities {
			//TODO: 可能存在账号还未创建的情况，需要考虑 UPSERT语义

			// 1. 更新余额
			err := g.db.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]any{
					"utime":   now,
					"balance": gorm.Expr("balance + ?", act.Account),
				}),
			}).Create(&Account{
				Uid:      act.Uid,
				Account:  act.Account,
				Type:     act.AccountType,
				Balance:  act.Account,
				Currency: act.Currency,
				Utime:    now,
				Ctime:    now,
			}).Error

			if err != nil {
				return err
			}

		}

		return tx.Create(&activities).Error

	})
}

// Account
// @Description: 虚拟账号表(后期该表会膨胀, 需要拆分)
type Account struct {
	Id      int64 `gorm:"primaryKey, autoIncrement"`
	Account int64 `gorm:"uniqueIndex:account_uid"`
	Type    uint8 `gorm:"uniqueIndex:account_uid"`
	Uid     int64 `gorm:"uniqueIndex:account_uid"`

	Balance  int64 //账户余额
	Currency string

	Utime int64
	Ctime int64
}

// AccountActivity
// @Description: 流水表
type AccountActivity struct {
	Id    int64  `gorm:"primaryKey, autoIncrement"`
	BizId int64  `gorm:"index:biz_type_id"`
	Biz   string `gorm:"index:biz_type_id"`
	Uid   int64
	// 正数表示入账 负数表示出账
	Account     int64 `gorm:"index:account_type"`
	AccountType uint8 `gorm:"index:account_type"`

	Amount   int64
	Currency string
	Utime    int64
	Ctime    int64
}
