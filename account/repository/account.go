package repository

import (
	"context"
	"kitbook/account/domain"
	"kitbook/account/repository/dao"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, credit domain.Credit) error
}

type CacheAccountRepository struct {
	dao dao.AccountDao
}

func NewCacheAccountRepository(dao dao.AccountDao) AccountRepository {
	return &CacheAccountRepository{
		dao: dao,
	}
}

func (c *CacheAccountRepository) AddCredit(ctx context.Context, credit domain.Credit) error {
	activities := make([]dao.AccountActivity, 0, len(credit.Items))
	for _, item := range credit.Items {
		activities = append(activities, c.ConvertsDaoAccountActivity(credit.Biz, credit.BizId, &item))
	}

	return c.dao.AddActivities(ctx, activities...)
}

func (c *CacheAccountRepository) ConvertsDomainCredit(act *dao.AccountActivity) domain.Credit {
	panic("todo")
}

func (c *CacheAccountRepository) ConvertsDaoAccountActivity(biz string, bizId int64, item *domain.CreditItem) dao.AccountActivity {
	return dao.AccountActivity{
		BizId:       bizId,
		Biz:         biz,
		Uid:         item.Uid,
		Account:     item.Account,
		AccountType: item.AccountType.AstoUint8(),
		Amount:      item.Amt,
		Currency:    item.Currency,
	}
}
