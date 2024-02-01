package wechat

import (
	"context"
	"kitbook/payment/domain"
)

type PaymentService interface {
	Prepay(ctx context.Context, pmt domain.Payment) (string, error)
}
