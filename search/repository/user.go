package repository

import (
	"context"
	"kitbook/search/domain"
	"kitbook/search/repository/dao"
)

type userSyncRepository struct {
	dao dao.UserDao
}

func NewUserSyncRepository(dao dao.UserDao) UserSyncRepository {
	return &userSyncRepository{
		dao: dao,
	}
}

func (u *userSyncRepository) InputUser(ctx context.Context, user domain.User) error {
	return u.dao.InputUser(ctx, u.convertsUserDao(&user))
}

func (u *userSyncRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userSyncRepository) convertsUserDao(user *domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		NikeName: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
	}
}
