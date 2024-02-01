package repository

import (
	"context"
	"kitbook/payment/domain"
	"kitbook/payment/repository/dao"
	"time"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, pmt domain.Payment) error
	UpdatePayment(ctx context.Context, pmt domain.Payment) error
	FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]domain.Payment, error)
}

type NativePaymentRepository struct {
	dao dao.PaymentDao
}

func NewNativePaymentRepository(dao dao.PaymentDao) PaymentRepository {
	return &NativePaymentRepository{
		dao: dao,
	}
}

func (n *NativePaymentRepository) CreatePayment(ctx context.Context, pmt domain.Payment) error {
	return n.dao.Insert(ctx, n.ConvertsPaymentDao(&pmt))
}

func (n *NativePaymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	return n.dao.Update(ctx, pmt.BizTradeNO, pmt.TxnID, pmt.Status.AstoUint8())
}

func (n *NativePaymentRepository) FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]domain.Payment, error) {
	pmtsDao, err := n.dao.FindExpiredPayment(ctx, beforeTime, offset, limit)
	if err != nil {
		return nil, err
	}

	pmts := make([]domain.Payment, len(pmtsDao))
	for i, p := range pmtsDao {
		pmts[i] = n.ConvertsPaymentDomain(&p)
	}
	return pmts, nil
}

func (n *NativePaymentRepository) ConvertsPaymentDomain(pmt *dao.Payment) domain.Payment {
	return domain.Payment{
		Amt: domain.Amount{
			Total:    pmt.Total,
			Currency: pmt.Currency,
		},
		BizTradeNO:  pmt.BizTradeNO,
		Description: pmt.Description,
		Status:      domain.PaymentStatus(pmt.Status),
		TxnID:       pmt.TxnID,
	}
}

func (n *NativePaymentRepository) ConvertsPaymentDao(pmt *domain.Payment) dao.Payment {
	return dao.Payment{
		Total:       pmt.Amt.Total,
		Currency:    pmt.Amt.Currency,
		BizTradeNO:  pmt.BizTradeNO,
		TxnID:       pmt.TxnID,
		Status:      pmt.Status.AstoUint8(),
		Description: pmt.Description,
	}
}
