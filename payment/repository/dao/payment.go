package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type PaymentDao interface {
	Insert(ctx context.Context, pmt Payment) error
	Update(ctx context.Context, bizTradeNo string, txnID string, status uint8) error
	FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]Payment, error)
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

// @func: FindExpiredPayment
// @date: 2024-02-02 03:01:48
// @brief: 批量查询超时过期订单
// @author: Kewin Li
// @receiver g
// @param ctx
// @param beforeTime
// @param offset
// @param limit
// @return Payment
// @return error
func (g *GormPaymentDao) FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]Payment, error) {
	var pmts []Payment
	err := g.db.WithContext(ctx).Where("utime < ?", beforeTime).
		Offset(offset).Limit(limit).Find(&pmts).Error
	return pmts, err
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
