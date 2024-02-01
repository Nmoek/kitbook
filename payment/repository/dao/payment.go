package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type PaymentDao interface {
	Insert(ctx context.Context, pmt Payment) error
	Update(ctx context.Context, bizTradeNo string, txnID string, status uint8) error
}

type GormPaymentDao struct {
	db *gorm.DB
}

func NewGormPaymentDao(db *gorm.DB) PaymentDao {
	return &GormPaymentDao{
		db: db,
	}
}

func (g *GormPaymentDao) Insert(ctx context.Context, pmt Payment) error {
	return g.db.WithContext(ctx).Create(&pmt).Error
}

func (g *GormPaymentDao) Update(ctx context.Context, bizTradeNo string, txnID string, status uint8) error {
	return g.db.WithContext(ctx).Model(&Payment{}).
		Where("biz_trade_no = ?", bizTradeNo).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": status,
			"txn_id": txnID,
		}).Error
}

// Payment
// @Description: 支付订单表
type Payment struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	BizTradeNO string `gorm:"uniqueIndex:tadeNo_uid_type_id"`
	UserId     int64  `gorm:"uniqueIndex:tadeNo_uid_type_id"`
	TxnID      string `gorm:"index"`

	Description string `gorm:"type:varchar(4096)"`
	Total       int64
	Currency    string
	Status      uint8

	Utime int64
	Ctime int64
}
