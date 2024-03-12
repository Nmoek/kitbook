package repository

import (
	"context"
	"kitbook/search/repository/dao"
)

type AnyRepository interface {
	Input(ctx context.Context, index string, docID string, data string) error
}

type anyRepository struct {
	dao dao.AnyDao
}

func NewAnyRepository(dao dao.AnyDao) AnyRepository {
	return &anyRepository{
		dao: dao,
	}
}

func (a *anyRepository) Input(ctx context.Context, index string, docID string, data string) error {
	return a.dao.Input(ctx, index, docID, data)
}
