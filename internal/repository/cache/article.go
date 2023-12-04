// Package cache
// @Description: 帖子查询缓存
package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/internal/domain"
	"time"
)

var (
	//go:embed lua/get_firstPage.lua
	luaGetFirstPage string
	//go:embed lua/set_firstPage.lua
	luaSetFirstPage string
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, userId int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}

//	@func: GetFirstPage
//	@date: 2023-12-04 22:09:07
//	@brief: 获取第一页数据缓存-非lua
//	@author: Kewin Li
//	@receiver r
//	@param ctx
//	@param userId
//	@return []domain.Article
//	@return error
//
// TODO: 实现lua脚本版本, 非lua脚本版本
func (r *RedisArticleCache) GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error) {

	// 1. 手动转换，多一层拷贝string-->[]byte
	//val, err := r.client.Get(ctx, r.createKey(userId)).Result()
	//if err != nil {
	//	return nil, err
	//}

	//2. unsafe转换, 内存零拷贝转换(推荐)
	val, err := r.client.Get(ctx, r.createKey(userId)).Bytes()

	var arts []domain.Article
	err = json.Unmarshal(val, &arts)

	return arts, err
}

// @func: GetFirstPageV1
// @date: 2023-12-04 22:09:51
// @brief: 获取第一页数据缓存-lua脚本
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return []domain.Article
// @return error
func (r *RedisArticleCache) GetFirstPageV1(ctx context.Context, userId int64) ([]domain.Article, error) {

	panic("implement me")

	key := r.createKey(userId)
	// TODO: lua脚本
	_, err := r.client.Eval(ctx, luaGetFirstPage, []string{key}).Result()
	return nil, err
}

// @func: SetFirstPage
// @date: 2023-12-04 22:09:32
// @brief: 设置第一页数据缓存-非lua
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error {

	// 生成文章摘要
	// 不需要缓存完整的Content
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].CreateAbstract()
	}

	val, err := json.Marshal(&arts)
	if err != nil {
		return err
	}

	// 注意! 超时时间的确定依赖于实际业务场景(流量大小、并发量高低)
	return r.client.Set(ctx, r.createKey(userId), val, 10*time.Minute).Err()

}

// @func: SetFirstPageV1
// @date: 2023-12-04 22:11:13
// @brief: 设置第一页数据缓存-lua脚本
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) SetFirstPageV1(ctx context.Context, userId int64) error {
	//TODO implement me
	panic("implement me")
}

// @func: DelFirstPage
// @date: 2023-12-04 22:38:35
// @brief: 清除第一页查询缓存数据-非lua
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) DelFirstPage(ctx context.Context, userId int64) error {

	return r.client.Del(ctx, r.createKey(userId)).Err()
}

// @func: DelFirstPageV1
// @date: 2023-12-04 22:38:56
// @brief: 清除第一页查询缓存数据-lua脚本
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) DelFirstPageV1(ctx context.Context, userId int64) error {
	panic("dont use it!")
}

// @func: createKey
// @date: 2023-12-04 22:37:41
// @brief: 生成帖子在Redis中的key
// @author: Kewin Li
// @receiver r
// @param userId
// @return string
func (r *RedisArticleCache) createKey(userId int64) string {
	return fmt.Sprintf("article:first_page:%d", userId)
}
