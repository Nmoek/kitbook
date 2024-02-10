package service

import (
	"context"
	"kitbook/account/domain"
)

type AccountService interface {
	Credit(ctx context.Context, credit domain.Credit) error
}
