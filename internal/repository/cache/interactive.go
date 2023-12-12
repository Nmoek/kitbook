package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const fieldReadCnt = "read_cnt"

type InteractiveCache interface {
	IncreaseReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

// @func: IncreaseReadCntIfPresent
// @date: 2023-12-11 23:30:02
// @brief: 阅读数+1
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param id
// @return error
func (r *RedisInteractiveCache) IncreaseReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	// 不太关注lua脚本的返回值, 就算字段不存在也会创建字段
	return r.client.Eval(ctx, luaIncrCnt, []string{r.createKey(biz, bizId)}, fieldReadCnt, 1).Err()
}

// @func: createKey
// @date: 2023-12-11 23:32:10
// @brief: 创建互动模块的key
// @author: Kewin Li
// @receiver r
// @param biz
// @param id
// @return string
func (r *RedisInteractiveCache) createKey(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
