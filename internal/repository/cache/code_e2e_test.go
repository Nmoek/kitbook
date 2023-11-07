// Package cache
// @Description: 单独测试lua脚本
package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// @func: TestRedisCodeCache_Set_e2e
// @date: 2023-11-06 00:45:19
// @brief: 单元测试-lua脚本
// @author: Kewin Li
// @param t
func TestRedisCodeCache_Set_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	testCases := []struct {
		name string

		before func(t *testing.T)

		after func(t *testing.T)

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "lua脚本设置成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:18762850585"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)

				leftTime, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, leftTime > time.Minute*9+time.Second*50)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "18762850585",
			code:  "123456",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			c := NewRedisCodeCache(rdb)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.NoError(t, err)
		})
	}
}
