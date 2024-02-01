package service

import (
	"context"
	"kitbook/payment/domain"
	"time"
)

type PaymentService interface {
	PrePay(ctx context.Context, pmt domain.Payment) (string, error)
	FindExpiredPayment(ctx context.Context, beforeTime time.Time, offset int, limit int) ([]domain.Payment, error)
}
