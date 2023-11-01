// Package cache
// @Description: 用户模块的缓存设计
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/internal/repository/dao"
	"time"
)

const ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (dao.User, error)
	Set(ctx context.Context, user dao.User) error
}

// RedisUserCache
// @Description: 基于Redis的缓存实现
type RedisUserCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

// @func: NewRedisUserCache
// @date: 2023-10-26 01:47:36
// @brief: 新建UserCache对象
// @author: Kewin Li
// @param cmd
// @return *RedisUserCache
func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: 15 * time.Minute,
	}
}

// @func: Get
// @date: 2023-10-26 01:55:15
// @brief: 缓存模块-查询用户信息
// @author: Kewin Li
// @receiver c
// @param id
// @return interface{}
// @return interface{}
func (c *RedisUserCache) Get(ctx context.Context, id int64) (dao.User, error) {
	k := c.Key(id)
	result := c.cmd.Get(ctx, k)
	if err := result.Err(); err != nil {
		return dao.User{}, err
	}

	// 反序列化
	var user dao.User
	err := json.Unmarshal([]byte(result.Val()), &user)

	return user, err
}

// @func: Set
// @date: 2023-10-26 02:00:27
// @brief: 缓存模块-插入用户信息
// @author: Kewin Li
// @receiver cache
// @param user
// @return error
func (c *RedisUserCache) Set(ctx context.Context, user dao.User) error {

	val, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	k := c.Key(user.Id)

	return c.cmd.Set(ctx, k, string(val), c.expiration).Err()
}

// @func: createKey
// @date: 2023-10-26 02:15:07
// @brief: 缓存模块-设计Key
// @author: Kewin Li
// @receiver cache
// @param id
// @return string
func (c *RedisUserCache) Key(id int64) string {
	//user-info-
	//user_info_
	//user.info.
	//user:info:
	return fmt.Sprintf("user_info_%d", id)
}
