package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/repository/cache/redismocks"
	"testing"
)

// @func: TestRedisCodeCache_Set
// @date: 2023-11-05 02:40:47
// @brief: 单元测试-redis插入缓存
// @author: Kewin Li
// @param t
func TestRedisCodeCache_Set(t *testing.T) {
	createKey := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "插入缓存成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				hdl := redis.NewCmd(context.Background())
				hdl.SetErr(nil)
				hdl.SetVal(int64(0))

				cmd.EXPECT().Eval(gomock.Any(),
					luaSetCode,
					[]string{createKey("test", "18762850585")}, "666666").Return(hdl)

				return cmd
			},

			biz:     "test",
			phone:   "18762850585",
			code:    "666666",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				hdl := redis.NewCmd(context.Background())
				hdl.SetErr(errors.New("redis错误"))

				cmd.EXPECT().Eval(gomock.Any(),
					luaSetCode,
					[]string{createKey("test", "18762850585")}, "666666").Return(hdl)

				return cmd
			},

			biz:     "test",
			phone:   "18762850585",
			code:    "666666",
			wantErr: errors.New("redis错误"),
		},
		{
			name: "验证码存在但没有设置过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				hdl := redis.NewCmd(context.Background())
				hdl.SetErr(nil)
				hdl.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(),
					luaSetCode,
					[]string{createKey("test", "18762850585")}, "666666").Return(hdl)

				return cmd
			},

			biz:     "test",
			phone:   "18762850585",
			code:    "666666",
			wantErr: errors.New("验证码没有设置过期时间"),
		},
		{
			name: "验证码发送频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				hdl := redis.NewCmd(context.Background())
				hdl.SetErr(nil)
				hdl.SetVal(int64(-2))
				cmd.EXPECT().Eval(gomock.Any(),
					luaSetCode,
					[]string{createKey("test", "18762850585")}, "666666").Return(hdl)

				return cmd
			},

			biz:     "test",
			phone:   "18762850585",
			code:    "666666",
			wantErr: ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cmd := tc.mock(ctrl)
			c := NewRedisCodeCache(cmd)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.Equal(t, tc.wantErr, err)

		})
	}
}
