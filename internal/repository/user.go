package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
)

var ErrDuplicateEmail = dao.ErrDuplicateEmail

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

// @func: Create
// @date: 2023-10-08 02:44:45
// @brief: 转发模块-数据创建
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param user
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}
