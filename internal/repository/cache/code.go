// Package cache
// @Description: 验证码缓存
package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

var ErrCodeSendTooMany = errors.New("验证码发送过于频繁")
var ErrCodeVerifyCntTooMany = errors.New("超过验证次数")
var ErrCodeNotRight = errors.New("验证码错误")

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// RedisCodeCache
// @Description: 基于Redis验证码实现
type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

// @func: Set
// @date: 2023-10-28 21:32:15
// @brief: 验证码-插入缓存
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @param code
// @return error
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	// 验证码存在但没有设置过期时间
	case -1:
		return errors.New("验证码没有设置过期时间")
	// 验证码发送过于频繁
	case -2:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

// @func: Verify
// @date: 2023-10-29 00:11:16
// @brief: 验证码-查询缓存
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @param code
// @return bool
// @return error
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		return false, errors.New("键值失效")
	// 验证码超出验证次数
	case -2:
		return false, ErrCodeVerifyCntTooMany
	// 验证码错误
	case -3:
		return false, ErrCodeNotRight
	default:
		return false, errors.New("未知错误")
	}
}

// @func: Key
// @date: 2023-10-28 21:32:57
// @brief: 生成验证码所对应的Key值
// @author: Kewin Li
// @receiver c
func (c *RedisCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
