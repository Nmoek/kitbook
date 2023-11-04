// Package service
// @Description: 用户服务-单元测试
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	repomocks "kitbook/internal/repository/mocks"
	"testing"
)

func TestNormalUserService_Login(t *testing.T) {
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
		{
			name: "登录成功",
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
		{
			name: "用户未找到",
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
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, errors.New("db错误"))
				return repo
			},

			email:    "123@qq.com",
			password: "Ljk741610",

			wantUser: domain.User{},
			wantErr:  errors.New("db错误"),
		},
		{
			name: "密码错误",
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
