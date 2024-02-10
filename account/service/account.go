package service

import (
	"context"
	"kitbook/account/domain"
	"kitbook/account/repository"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewaccountService(repo repository.AccountRepository) AccountService {
	return &accountService{
		repo: repo,
	}
}

func (a *accountService) Credit(ctx context.Context, credit domain.Credit) error {
	return a.repo.AddCredit(ctx, credit)
}
