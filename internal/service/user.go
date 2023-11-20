package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
)

var (
	ErrDuplicateUser         = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户名或密码不正确") //
	ErrInvalidUserAccess     = errors.New("非法用户访问")
)

// UserService
// @Description: 普通用户、超级用户、VIP用户
type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email string, passwaord string) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	SignupOrLoginWithPhone(ctx context.Context, phone string) (domain.User, error)
	SignupOrLoginWithWechat(ctx context.Context, info domain.WechtInfo) (domain.User, error)
}

// NormalUserService
// @Description: 普通用户实现
type NormalUserService struct {
	repo repository.UserRepository //一个服务只会有一个repository

}

// @func: NewNormalUserService
// @date: 2023-10-08 02:42:59
// @brief: 创建新的服务对象
// @author: Kewin Li
// @param repo
// @return *NormalUserService
func NewNormalUserService(repo repository.UserRepository) UserService {
	return &NormalUserService{
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
func (svc *NormalUserService) Signup(ctx context.Context, user domain.User) error {
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
func (svc *NormalUserService) Login(ctx context.Context, email string, passwaord string) (domain.User, error) {
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
	if err == bcrypt.ErrMismatchedHashAndPassword {
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
func (svc *NormalUserService) Edit(ctx context.Context, user domain.User) error {
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
func (svc *NormalUserService) Profile(ctx context.Context, id int64) (domain.User, error) {

	return svc.repo.FindByID(ctx, id)

}

// @func: SignupOrLoginWithPhone
// @date: 2023-10-29 22:43:51
// @brief: 通过手机号注册/登录
// @author: Kewin Li
// @receiver svc
// @param ctx
// @param phone
func (svc *NormalUserService) SignupOrLoginWithPhone(ctx context.Context, phone string) (domain.User, error) {

	// 默认认为大多数用户式已经注册的
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err == nil || err != repository.ErrUserNotFound {
		return user, err
	}

	// 用户不存在
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})

	// err有两种可能:
	// 1. 唯一索引冲突, 注册用户失败
	// 2. 其他错误

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	// 如果是已存在用户，查找出来即可
	// 如果是未注册用户，可能存在主从库延迟问题，一般主库写入，从库查出
	//TODO: 主从延迟问题，强制本次查询走主库
	return svc.repo.FindByPhone(ctx, phone)
}

// @func: SignupOrLoginWithWechat
// @date: 2023-11-12 01:32:22
// @brief: 通过微信号注册/ 登录
// @author: Kewin Li
// @receiver n
// @param ctx
// @param info
// @return domain.User
// @return error
func (svc *NormalUserService) SignupOrLoginWithWechat(ctx context.Context, info domain.WechtInfo) (domain.User, error) {
	// 默认认为大多数用户式已经注册的
	user, err := svc.repo.FindByWechat(ctx, info.Openid)
	if err == nil || err != repository.ErrUserNotFound {
		return user, err
	}

	// 用户不存在
	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: info,
	})

	// err有两种可能:
	// 1. 唯一索引冲突, 注册用户失败
	// 2. 其他错误

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}

	// 如果是已存在用户，查找出来即可
	// 如果是未注册用户，可能存在主从库延迟问题，一般主库写入，从库查出
	//TODO: 主从延迟问题，强制本次查询走主库
	return svc.repo.FindByWechat(ctx, info.Openid)
}
