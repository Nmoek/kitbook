// Package repository
// @Description: 数据转发-用户模块-单元测试
package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
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

// @func: TestCacheUserRepository_FindByID
// @date: 2023-11-04 21:39:39
// @brief: 单元测试-通过ID搜记录
// @author: Kewin Li
// @param t
func TestCacheUserRepository_FindByID(t *testing.T) {
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
		{
			name: "查询成功,缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(testUser, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(testUser, nil)

				c.EXPECT().Set(gomock.Any(), testUser).Return(nil)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr:  nil,
			wantUser: ConvertsDomainUser(&testUser),
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存命中
				c.EXPECT().Get(gomock.Any(), id).Return(testUser, nil)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr:  nil,
			wantUser: ConvertsDomainUser(&testUser),
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(testUser, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(dao.User{}, gorm.ErrRecordNotFound)

				return d, c

			},

			id:  1,
			ctx: context.Background(),

			wantErr: ErrUserNotFound,
		},
		{
			name: "回写缓存失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {

				d := daomocks.NewMockUserDao(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)

				// 缓存未命中
				c.EXPECT().Get(gomock.Any(), id).Return(testUser, cache.ErrKeyNotExist)

				d.EXPECT().FindByID(gomock.Any(), id).Return(testUser, nil)

				c.EXPECT().Set(gomock.Any(), testUser).Return(errors.New("redis出错"))
				return d, c

			},

			id:       1,
			ctx:      context.Background(),
			wantUser: ConvertsDomainUser(&testUser),
			wantErr:  nil,
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
