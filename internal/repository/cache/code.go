// Package cache
// @Description: 验证码缓存
package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"strconv"
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
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{createKey(biz, phone)}, code).Int()
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
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{createKey(biz, phone)}, code).Int()
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

// LocalCodeCache
// @Description: 基于本地内存缓存实现
type LocalCodeCache struct {
	mem    *freecache.Cache
	expire int
}

func NewLocalCodeCache(mem *freecache.Cache) CodeCache {
	return &LocalCodeCache{
		mem:    mem,
		expire: 600, // 10min
	}
}

// @func: Set
// @date: 2023-11-01 23:39:59
// @brief: 本地缓存-设置验证码缓存
// @author: Kewin Li
// @receiver l
// @param ctx
// @param biz
// @param phone
// @param code
// @return error
func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 先检查key是否存在
	// 如果存在, 检查过期时间
	// 如果不存在, 直接设置
	key := createKey(biz, phone)
	_, err := l.mem.Get([]byte(key))
	if err != nil && err != freecache.ErrNotFound {
		// TODO: 日志埋点
		return err
	}

	leftTime, err := l.mem.TTL([]byte(key))
	if err != nil && err != freecache.ErrNotFound {
		// TODO: 日志埋点
		return err
	}
	// 还未过1min
	if leftTime > 540 {
		// 返回发送频繁
		return ErrCodeSendTooMany
	}

	// 缓存验证码
	err = l.mem.Set([]byte(key), []byte(code), l.expire)
	if err != nil {
		return err
	}

	// 缓存验证码的验证上限
	key += ":cnt"
	val := strconv.FormatInt(3, 10)

	return l.mem.Set([]byte(key), []byte(val), l.expire)
}

// @func: Verify
// @date: 2023-11-01 23:40:08
// @brief: 本地缓存-校验验证码缓存
// @author: Kewin Li
// @receiver l
// @param ctx
// @param biz
// @param phone
// @param code
// @return bool
// @return error
func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 先检查key是否存在
	// 不存在, 返回错误
	// 如果key存在
	// 取出对应val判断是否正确
	// 正确，删除key
	// 不正常, 验证次数减1, 并判断是否超出验证限制

	key := createKey(biz, phone)
	val, err := l.mem.Get([]byte(key))
	if err != nil {
		return false, err
	}

	if string(val) != code {
		key += ":cnt"

		cntStr, leftTime, err := l.mem.GetWithExpiration([]byte(key))
		if err != nil {
			return false, err
		}

		cnt, _ := strconv.ParseInt(string(cntStr), 10, 64)
		if cnt <= 0 {
			return false, ErrCodeVerifyCntTooMany
		}

		cntStr = []byte(strconv.FormatInt(cnt-1, 10))

		err = l.mem.Set([]byte(key), cntStr, int(leftTime))
		if err != nil {
			return false, err
		}

		return false, ErrCodeNotRight
	}

	// 验证成功需要删除键值对
	if l.mem.Del([]byte(key)) {
		if l.mem.Del([]byte(key + ":cnt")) {
			return true, nil
		}
	}

	// TODO: 日志埋点，key删除失败

	//删除失败将验证次数置0, 过期时间置为1s
	err = l.mem.Set([]byte(key+":cnt"), []byte("0"), 1)
	if err != nil {
		return false, err
	}

	return true, nil
}

// @func: createKey
// @date: 2023-10-28 21:32:57
// @brief: 生成验证码所对应的Key值
// @author: Kewin Li
// @receiver c
func createKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
