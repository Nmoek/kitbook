// Package repository
// @Description: 数据转发-用户模块-单元测试
package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
	cachemocks "kitbook/internal/repository/cache/mocks"
	"kitbook/internal/repository/dao"
	daomocks "kitbook/internal/repository/dao/mocks"
	"testing"
	"time"
)

type UserRepositorySuite struct {
	suite.Suite
}

// @func: TestCreate
// @date: 2024-01-19 00:06:34
// @brief: 单元测试-repository层-创建用户
// @author: Kewin Li
// @receiver u
func (u *UserRepositorySuite) TestCreate() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		user    domain.User
		wantErr error
	}{
		// 创建用户成功
		{
			name: "User created successfully",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				return d, c
			},

			user: domain.User{
				Email:    "1@qq.com",
				Password: "12345",
			},
		},
		// 用户已存在
		{
			name: "User already exists",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(dao.ErrDuplicateUser)
				return d, c
			},

			wantErr: dao.ErrDuplicateUser,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("数据库错误"))
				return d, c
			},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCacheUserRepository(d, c)
			err := repo.Create(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

// @func: TestFindByEmail
// @date: 2024-01-19 00:59:00
// @brief: 单元测试-repository层-按邮箱查询
// @author: Kewin Li
// @receiver u
func (u *UserRepositorySuite) TestFindByEmail() {
	t := u.T()
	now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		email string

		wantUser domain.User
		wantErr  error
	}{
		// 查询成功
		{
			name: "Enquiry Successful",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(dao.User{
						Id: 1,
						Email: sql.NullString{
							String: "1@qq.com",
							Valid:  true,
						},
						Password: "12345",
						Birthday: now,
					}, nil)

				return d, c
			},
			wantUser: domain.User{
				Id:       1,
				Email:    "1@qq.com",
				Password: "12345",
				Birthday: time.UnixMilli(now),
			},
		},
		// 用户不存在
		{
			name: "User not Find",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(dao.User{}, ErrUserNotFound)

				return d, c
			},

			wantErr: ErrUserNotFound,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				d.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).
					Return(dao.User{}, errors.New("数据库错误"))

				return d, c
			},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCacheUserRepository(d, c)
			user, err := repo.FindByEmail(context.Background(), tc.email)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}
}

// @func: TestFindByID
// @date: 2024-01-19 00:05:49
// @brief: 单元测试-repository层-按userId查询
// @author: Kewin Li
// @receiver u
func (u *UserRepositorySuite) TestFindByID() {
	t := u.T()
	nowMs := time.Now().UnixMilli()
	id := int64(1)
	testUser := dao.User{
		Id:       id,
		Email:    sql.NullString{String: "123@qq.com", Valid: true},
		Phone:    sql.NullString{String: "1851561615", Valid: true},
		Password: "Ljk741610",
		Birthday: nowMs,
		AboutMe:  "这是一个mock测试",
		Ctime:    nowMs,
		Utime:    nowMs,
	}

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		id  int64
		ctx context.Context

		wantErr  error
		wantUser domain.User
	}{
		// 缓存未命中, 数据库查询成功, 缓存回写成功
		{
			name: "Cache not hit, database query successfully",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				c.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(testUser, nil)

				c.EXPECT().Set(gomock.Any(), ConvertsDomainUser(&testUser)).Return(nil)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr:  nil,
			wantUser: ConvertsDomainUser(&testUser),
		},
		// 缓存命中
		{
			name: "Cache hit",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存命中
				c.EXPECT().Get(gomock.Any(), id).Return(ConvertsDomainUser(&testUser), nil)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr:  nil,
			wantUser: ConvertsDomainUser(&testUser),
		},
		// 缓存未命中, 数据库未查询到
		{
			name: "Cache not hit, database not queried",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(dao.User{}, gorm.ErrRecordNotFound)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr: ErrUserNotFound,
		},
		// 缓存未命中, 数据库发生其他错误
		{
			name: "Cache not hit, other database errors occurred",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(dao.User{}, errors.New("数据库错误"))

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr: errors.New("数据库错误"),
		},
		// 缓存回写失败
		{
			name: "Cache Write Back Failure",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(testUser, nil)

				c.EXPECT().Set(gomock.Any(), gomock.Any()).Return(errors.New("redis出错"))
				return d, c

			},

			id:       1,
			ctx:      context.Background(),
			wantUser: ConvertsDomainUser(&testUser),
			wantErr:  nil,
		},
		// 缓存发生其他错误
		{
			name: "Other errors occurred in the cache",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(domain.User{}, errors.New("缓存发生其他错误"))

				return d, c

			},

			id:      1,
			ctx:     context.Background(),
			wantErr: errors.New("缓存发生其他错误"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCacheUserRepository(d, c)
			user, err := repo.FindById(tc.ctx, tc.id)
			//
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)

		})
	}

}

// @func: FindByPhone
// @date: 2024-01-19 01:59:10
// @brief: 单元测试-repository层-按手机号查询
// @author: Kewin Li
// @receiver u
func (u *UserRepositorySuite) TestFindByPhone() {
	t := u.T()

	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		phone string

		wantErr  error
		wantUser domain.User
	}{
		// 查询成功
		{
			name: "Enquiry Successful",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).Return(dao.User{
					Id: 1,
					Phone: sql.NullString{
						String: "111",
					},
				}, nil)

				return d, c

			},

			phone: "111",
			wantUser: domain.User{
				Id:       1,
				Phone:    "111",
				Birthday: time.UnixMilli(0),
			},
		},
		// 查询失败
		{
			name: "Enquiry Fail",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByPhone(gomock.Any(), gomock.Any()).Return(dao.User{}, errors.New("数据库错误"))

				return d, c

			},

			phone:   "222",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCacheUserRepository(d, c)
			user, err := repo.FindByPhone(context.Background(), tc.phone)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

// @func: TestFindByWechat
// @date: 2024-01-19 01:59:10
// @brief: 单元测试-repository层-按微信号查询
// @author: Kewin Li
// @receiver u
func (u *UserRepositorySuite) TestFindByWechat() {
	t := u.T()

	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)

		openid string

		wantErr  error
		wantUser domain.User
	}{
		// 查询成功
		{
			name: "Enquiry Successful",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).Return(dao.User{
					Id: 1,
					Openid: sql.NullString{
						String: "111",
					},
				}, nil)

				return d, c

			},

			openid: "111",
			wantUser: domain.User{
				Id: 1,
				WechatInfo: domain.WechatInfo{
					Openid: "111",
				},
				Birthday: time.UnixMilli(0),
			},
		},
		// 查询失败
		{
			name: "Enquiry Fail",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				d.EXPECT().FindByWechat(gomock.Any(), gomock.Any()).Return(dao.User{}, errors.New("数据库错误"))

				return d, c

			},

			openid:  "222",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			d, c := tc.mock(ctrl)
			repo := NewCacheUserRepository(d, c)
			user, err := repo.FindByWechat(context.Background(), tc.openid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, &UserRepositorySuite{})
}
