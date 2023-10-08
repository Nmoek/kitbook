package service

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
)

var ErrDuplicateEmail = repository.ErrDuplicateEmail

type UserService struct {
	repo *repository.UserRepository //一个服务只会有一个repository
}

// @func: NewUserService
// @date: 2023-10-08 02:42:59
// @brief: 创建新的服务对象
// @author: Kewin Li
// @param repo
// @return *UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// @func: Signup
// @date: 2023-10-08 02:39:12
// @brief: 用户注册服务
// @author: Kewin Li
// @receiver svc
// @param ctx  TODO: 这个参数是做干嘛的？？？
// @param user
// @return error
func (svc *UserService) Signup(ctx context.Context, user domain.User) error {
	// 加密服务
	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Printf("crypt password %s \n", cryptPassword)
	user.Password = string(cryptPassword)
	return svc.repo.Create(ctx, user)
}
