package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户名或密码不正确") //
	ErrInvalidUserAccess     = errors.New("非法用户访问")
)

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
// @param ctx  TODO: 这个参数是做干嘛的？
// @param user
// @return error
func (svc *UserService) Signup(ctx context.Context, user domain.User) error {
	// 加密服务
	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(cryptPassword)
	return svc.repo.Create(ctx, user)

}

// @func: Login
// @date: 2023-10-09 01:49:35
// @brief: 用户登录服务
// @author: Kewin Li
// @receiver svc
// @param ctx
// @param user
// @return error
func (svc *UserService) Login(ctx context.Context, email string, passwaord string) (domain.User, error) {
	// 1. 根据邮箱找到用户
	findUser, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	// 2. 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(passwaord))
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	return findUser, nil
}

// @func: Edit
// @date: 2023-10-12 03:44:44
// @brief: 修改个人信息服务
// @author: Kewin Li
// @receiver svc
// @param ctx
// @param info
func (svc *UserService) Edit(ctx context.Context, user domain.User) error {
	return svc.repo.UpdatePersonalInfo(ctx, user)
}

// @func: Profile
// @date: 2023-10-14 17:58:53
// @brief: 查看个人信息服务
// @author: Kewin Li
// @receiver svc
// @param ctx
// @param id
// @return interface{}
// @return interface{}
func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {

	return svc.repo.FindByID(ctx, id)

}

// @func: CheckID
// @date: 2023-10-12 04:03:56
// @brief: 查询拿到的ID是否存在
// @author: Kewin Li
// @receiver svc
// @param ctx
// @param id
func (svc *UserService) CheckID(ctx context.Context, id int64) error {
	_, err := svc.repo.FindByID(ctx, id)
	return err
}
