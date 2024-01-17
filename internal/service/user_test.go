// Package service
// @Description: 用户服务-单元测试
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	repomocks "kitbook/internal/repository/mocks"
	"testing"
)

type UserServiceSuite struct {
	suite.Suite
}

// @func: TestSignup
// @date: 2024-01-18 03:45:16
// @brief: 单元测试-service层-注册
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestSignup() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		user domain.User

		wantErr error
	}{
		// 插入注册信息成功
		{
			name: "Insert registration information successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				return repo
			},

			user: domain.User{
				Email:    "1@qq.com",
				Password: "123456",
			},
		},
		// 用户重复注册
		{
			name: "Duplicate user registrations",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repository.ErrDuplicateUser)
				return repo
			},

			user: domain.User{
				Email:    "2@qq.com",
				Password: "123456",
			},
			wantErr: repository.ErrDuplicateUser,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)

				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("数据其他未知错误"))
				return repo
			},

			user: domain.User{
				Email:    "3@qq.com",
				Password: "123456",
			},
			wantErr: errors.New("数据其他未知错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)

			svc := NewNormalUserService(repo)
			err := svc.Signup(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

// @func: TestLogin
// @date: 2024-01-18 03:43:31
// @brief: 单元测试-service层-登录
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestLogin() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(strl *gomock.Controller) repository.UserRepository

		// 预期输入
		ctx      context.Context
		email    string
		password string

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		// 用户登录成功
		{
			name: "User Login Success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 加密后正确的密码
					Password: "$2a$10$yl7gU376QU/DYLL/zgBOtuPA3eQWn9qIK6LFp9jg797iUfCcxUQ3i",
					Phone:    "123151616",
				}, nil)

				return repo
			},

			email:    "123@qq.com",
			password: "Ljk741610",

			wantUser: domain.User{
				Email: "123@qq.com",
				// 加密后正确的密码
				Password: "$2a$10$yl7gU376QU/DYLL/zgBOtuPA3eQWn9qIK6LFp9jg797iUfCcxUQ3i",
				Phone:    "123151616",
			},
			wantErr: nil,
		},
		// 用户未找到
		{
			name: "User not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},

			email:    "123@qq.com",
			password: "Ljk741610",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		// 用户密码校验不通过
		{
			name: "User password verification failed",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 加密后正确的密码
					Password: "$2a$10$yl7gU376QU/DYLL/zgBOtuPA3eQWn9qIK6LFp9jg797iUfCcxUQ3i",
					Phone:    "123151616",
				}, ErrInvalidUserOrPassword)
				return repo
			},

			email:    "123@qq.com",
			password: "Ljk74161031133",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},

			email:    "123@qq.com",
			password: "Ljk741610",

			wantUser: domain.User{},
			wantErr:  errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewNormalUserService(repo)

			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}

// @func: TestEdit
// @date: 2024-01-18 04:10:14
// @brief: 单元测试-service层-修改个人信息
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestEdit() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(strl *gomock.Controller) repository.UserRepository

		// 预期输入
		user domain.User

		// 预期输出
		wantErr error
	}{
		// 修改信息成功
		{
			name: "Modify information successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().UpdatePersonalInfo(gomock.Any(), gomock.Any()).Return(nil)

				return repo
			},
			user: domain.User{
				Email:    "1@qq.com",
				Password: "123456",
			},
		},
		// 用户未找到
		{
			name: "User not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().UpdatePersonalInfo(gomock.Any(), gomock.Any()).Return(repository.ErrUserNotFound)
				return repo
			},

			user: domain.User{
				Email:    "2@qq.com",
				Password: "123456",
			},

			wantErr: repository.ErrUserNotFound,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().UpdatePersonalInfo(gomock.Any(), gomock.Any()).Return(errors.New("数据库错误"))
				return repo
			},

			user: domain.User{
				Email:    "3@qq.com",
				Password: "123456",
			},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewNormalUserService(repo)

			err := svc.Edit(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}

// @func: TestProfile
// @date: 2024-01-18 04:18:33
// @brief: 单元测试-service层-查看个人信息
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestProfile() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(strl *gomock.Controller) repository.UserRepository

		id int64

		wantUser domain.User
		wantErr  error
	}{
		// 查询信息成功
		{
			name: "Modify information successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Email:    "1@qq.com",
						Password: "123456",
					}, nil)

				return repo
			},
			id: 1,
			wantUser: domain.User{
				Email:    "1@qq.com",
				Password: "123456",
			},
		},
		// 用户未找到
		{
			name: "User not found",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},

			id:      2,
			wantErr: ErrInvalidUserAccess,
		},
		// 缓存错误
		{
			name: "Cache errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), gomock.Any()).
					Return(domain.User{}, nil)
				return repo
			},
			id: 3,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindById(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},
			id: 4,

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewNormalUserService(repo)

			user, err := svc.Profile(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}

// @func: TestSignupOrLoginWithPhone
// @date: 2024-01-18 04:26:23
// @brief: 单元测试-service层-通过手机号注册/登录
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestSignupOrLoginWithPhone() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(strl *gomock.Controller) repository.UserRepository

		phone string

		wantUser domain.User
		wantErr  error
	}{
		// 通过手机号登录成功
		{
			name: "login by phone Successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Phone:    "111",
						Password: "123456",
					}, nil)

				return repo
			},
			phone: "111",
			wantUser: domain.User{
				Phone:    "111",
				Password: "123456",
			},
		},
		// 用户未注册, 通过手机号注册成功
		{
			name: "User is not registered, registration by phone number successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 1. 第一次没查到
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)

				// 2. 创建
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				// 3. 第二次查到
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Phone: "222",
					}, nil)

				return repo
			},

			phone: "222",
			wantUser: domain.User{
				Phone: "222",
			},
		},
		// 用户重复注册（并发问题会出现）
		{
			name: "Duplicate user registrations",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 1. 第一次没查到
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)

				// 2. 创建也不成功失败, 注意:
				// 此时说明有节点创建成功了, 第一次因为脏数据没读到
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(repository.ErrDuplicateUser)

				// 3. 第二次查成功
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Phone: "333",
					}, nil)

				return repo
			},
			phone: "333",
			wantUser: domain.User{
				Phone: "333",
			},
		},
		// 第一次查询伴随其他错误
		{
			name: "First query with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("数据库错误"))

				return repo
			},
			phone:   "555",
			wantErr: errors.New("数据库错误"),
		},
		// 第一次查询没查到, 创建伴随其他错误
		{
			name: "First query didn't find it, created with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("数据库错误"))

				return repo
			},
			phone:   "666",
			wantErr: errors.New("数据库错误"),
		},
		// 第二次查询伴随其他错误
		{
			name: "Second query with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
				repo.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},
			phone:   "777",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewNormalUserService(repo)

			user, err := svc.SignupOrLoginWithPhone(context.Background(), tc.phone)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}

// @func: TestSignupOrLoginWithPhone
// @date: 2024-01-18 04:26:23
// @brief: 单元测试-service层-通过微信号注册/登录
// @author: Kewin Li
// @receiver u
func (u *UserServiceSuite) TestSignupOrLoginWithWechat() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(strl *gomock.Controller) repository.UserRepository

		info domain.WechatInfo

		wantUser domain.User
		wantErr  error
	}{
		// 通过微信号登录成功
		{
			name: "login by wechat Successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{
						WechatInfo: domain.WechatInfo{
							Openid: "111",
						},
					}, nil)

				return repo
			},
			info: domain.WechatInfo{
				Openid: "111",
			},
			wantUser: domain.User{
				WechatInfo: domain.WechatInfo{
					Openid: "111",
				},
			},
		},
		// 用户未注册, 通过微信号注册成功
		{
			name: "User is not registered, registration by wechat number successfully",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 1. 第一次没查到
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)

				// 2. 创建
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				// 3. 第二次查到
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{
						WechatInfo: domain.WechatInfo{
							Openid: "222",
						},
					}, nil)

				return repo
			},

			info: domain.WechatInfo{
				Openid: "222",
			},
			wantUser: domain.User{
				WechatInfo: domain.WechatInfo{
					Openid: "222",
				},
			},
		},
		// 用户重复注册（并发问题会出现）
		{
			name: "Duplicate user registrations",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 1. 第一次没查到
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)

				// 2. 创建也不成功失败, 注意:
				// 此时说明有节点创建成功了, 第一次因为脏数据没读到
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(repository.ErrDuplicateUser)

				// 3. 第二次查成功
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{
						WechatInfo: domain.WechatInfo{
							Openid: "333",
						},
					}, nil)

				return repo
			},
			info: domain.WechatInfo{
				Openid: "333",
			},
			wantUser: domain.User{
				WechatInfo: domain.WechatInfo{
					Openid: "333",
				},
			},
		},
		// 第一次查询伴随其他错误
		{
			name: "First query with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("数据库错误"))

				return repo
			},

			wantErr: errors.New("数据库错误"),
		},
		// 第一次查询没查到, 创建伴随其他错误
		{
			name: "First query didn't find it, created with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(errors.New("数据库错误"))

				return repo
			},

			wantErr: errors.New("数据库错误"),
		},
		// 第二次查询伴随其他错误
		{
			name: "Second query with other errors",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, repository.ErrUserNotFound)
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil)
				repo.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewNormalUserService(repo)

			user, err := svc.SignupOrLoginWithWechat(context.Background(), tc.info)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}

func TestUserService(t *testing.T) {
	suite.Run(t, &UserServiceSuite{})
}

// @func: TestPasswordEncrypt
// @date: 2023-11-04 19:34:06
// @brief: 单元测试-密码加密
// @author: Kewin Li
// @param t
func TestPasswordEncrypt(t *testing.T) {
	password := []byte("Ljk741610")

	encodePassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)

	fmt.Printf("%v \n", string(encodePassword))

	err = bcrypt.CompareHashAndPassword(encodePassword, []byte("Ljk741610"))
	assert.NoError(t, err)
}
