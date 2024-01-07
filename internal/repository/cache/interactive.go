package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/internal/domain"
	"strconv"
	"time"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
)

type InteractiveCache interface {
	IncreaseReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectionCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrCollectionCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
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

// @func: IncreaseLikeCntIfPresent
// @date: 2023-12-13 22:11:39
// @brief: 点赞数+1
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (r *RedisInteractiveCache) IncreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt, []string{r.createKey(biz, bizId)}, fieldLikeCnt, 1).Err()
}

// @func: DecreaseLikeCntIfPresent
// @date: 2023-12-13 22:11:47
// @brief: 点赞数-1
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (r *RedisInteractiveCache) DecreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt, []string{r.createKey(biz, bizId)}, fieldLikeCnt, -1).Err()

}

// @func: IncrCollectionCntIfPresent
// @date: 2023-12-14 02:04:49
// @brief: 收藏数+1
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (r *RedisInteractiveCache) IncrCollectionCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt, []string{r.createKey(biz, bizId)}, fieldCollectCnt, 1).Err()
}

// @func: DecrCollectionCntIfPresent
// @date: 2023-12-14 02:04:57
// @brief: 收藏数-1
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param bizId
// @param userId
// @return error
func (r *RedisInteractiveCache) DecrCollectionCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt, []string{r.createKey(biz, bizId)}, fieldCollectCnt, -1).Err()

}

// @func: Get
// @date: 2023-12-15 17:05:05
// @brief: 获取缓存互动模块数据
// @author: Kewin Li
// @receiver r
// @param ctx
// @param biz
// @param bizId
// @return dao.Interactive
// @return error
func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	key := r.createKey(biz, bizId)

	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	// 数据不存在的情况
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}

	return r.convertsDomainIntrFromCache(res)
}

// @func: Set
// @date: 2023-12-15 17:41:14
// @brief: 写入缓存互动模块数据
// @author: Kewin Li
// @receiver r
// @param ctx
// @param intr
// @return error
func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := r.createKey(biz, bizId)
	err := r.client.HSet(ctx, key,
		bizId, intr.BizId,
		fieldReadCnt, intr.ReadCnt,
		fieldLikeCnt, intr.LikeCnt,
		fieldCollectCnt, intr.CollectCnt).Err()

	if err != nil {
		return err
	}

	// 设置过期时间
	return r.client.Expire(ctx, key, 15*time.Minute).Err()
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

// @func: convertsDomainIntrFromCache
// @date: 2023-12-15 17:52:01
// @brief: Interactive转换 redis cache ---> domain
// @author: Kewin Li
// @receiver r
// @param res
// @return domain.Interactive
func (r *RedisInteractiveCache) convertsDomainIntrFromCache(res map[string]string) (domain.Interactive, error) {
	bizId, err := strconv.ParseInt(res["biz_id"], 10, 64)
	readCnt, err := strconv.ParseInt(res[fieldReadCnt], 10, 64)
	likeCnt, err := strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	collectCnt, err := strconv.ParseInt(res[fieldCollectCnt], 10, 64)

	return domain.Interactive{
		BizId:      bizId,
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, err
}
