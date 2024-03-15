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

const ErrKeyNotExist = redis.Nil

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
	GetById(ctx context.Context, artId int64) (domain.Article, error)
	SetById(ctx context.Context, art domain.Article) error
	GetPubById(ctx context.Context, artId int64) (domain.Article, error)
	SetPubById(ctx context.Context, art domain.Article) error
	SetPub(ctx context.Context, art domain.Article) error
	DelPub(ctx context.Context, id int64) error
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
	//val, err := r.client.Get(ctx, r.createFirstPageKey(userId)).Result()
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal([]byte(val), &arts)

	//2. unsafe转换, 内存零拷贝转换(推荐)
	val, err := r.client.Get(ctx, r.createFirstPageKey(userId)).Bytes()

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

	key := r.createFirstPageKey(userId)
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
	return r.client.Set(ctx, r.createFirstPageKey(userId), val, 10*time.Minute).Err()

}

// @func: SetFirstPageV1
// @date: 2023-12-04 22:11:13
// @brief: 设置第一页列表数据缓存-lua脚本
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
// @brief: 清除第一页列表查询缓存数据-非lua
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) DelFirstPage(ctx context.Context, userId int64) error {

	return r.client.Del(ctx, r.createFirstPageKey(userId)).Err()
}

// @func: DelFirstPageV1
// @date: 2023-12-04 22:38:56
// @brief: 清除第一页列表查询缓存数据-lua脚本
// @author: Kewin Li
// @receiver r
// @param ctx
// @param userId
// @return error
func (r *RedisArticleCache) DelFirstPageV1(ctx context.Context, userId int64) error {
	panic("dont use it!")
}

// @func: GetById
// @date: 2023-12-05 23:00:28
// @brief: 获取列表详情缓存-按Id
// @author: Kewin Li
// @receiver r
// @param ctx
// @param artId
// @return Article
// @return error
func (r *RedisArticleCache) GetById(ctx context.Context, artId int64) (domain.Article, error) {
	var art domain.Article
	// 1. 手动转换
	//val, err := r.client.Get(ctx, r.createPreCacheKey(artId)).Result()
	//if err != nil {
	//	return art, err
	//}
	//err = json.Unmarshal([]byte(val), &art)

	// 2. unsafe 转换 少一层拷贝
	val, err := r.client.Get(ctx, r.createPreCacheKey(artId)).Bytes()
	if err != nil {
		return art, err
	}

	err = json.Unmarshal(val, &art)

	return art, err
}

// @func: SetById
// @date: 2023-12-05 23:00:23
// @brief: 设置列表详情缓存-按Id
// @author: Kewin Li
// @receiver r
// @param ctx
// @param art
// @return error
func (r *RedisArticleCache) SetById(ctx context.Context, art domain.Article) error {

	val, err := json.Marshal(&art)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createPreCacheKey(art.Id), val, time.Minute).Err()
}

// @func: GetPubById
// @date: 2023-12-06 22:59:52
// @brief: 帖子查询-获取读者帖子缓存
// @author: Kewin Li
// @receiver r
// @param ctx
// @param artId
// @return domain.Article
// @return error
func (r *RedisArticleCache) GetPubById(ctx context.Context, artId int64) (domain.Article, error) {
	var art domain.Article

	val, err := r.client.Get(ctx, r.createPreCacheKey(artId)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}

	// key不存在
	if len(val) == 0 {
		return domain.Article{}, ErrKeyNotExist
	}

	err = json.Unmarshal(val, &art)

	return art, err
}

// @func: SetPubById
// @date: 2023-12-06 23:00:14
// @brief: 帖子查询-设置读者帖子缓存
// @author: Kewin Li
// @receiver r
// @param ctx
// @param art
// @return error
func (r *RedisArticleCache) SetPubById(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(&art)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createPreCacheKey(art.Id), val, 10*time.Minute).Err()
}

// @func: SetPub
// @date: 2023-12-07 02:39:48
// @brief: 帖子查询-发表时进行缓存
// @author: Kewin Li
// @receiver r
// @param ctx
// @param art
// @return error
func (r *RedisArticleCache) SetPub(ctx context.Context, art domain.Article) error {
	val, err := json.Marshal(&art)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createPreCacheKey(art.Id), val, 10*time.Minute).Err()
}

// @func: DelPub
// @date: 2024-03-15 17:17:30
// @brief: 删除缓存
// @author: Kewin Li
// @receiver r
// @param ctx
// @param art
// @return error
func (r *RedisArticleCache) DelPub(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.createPreCacheKey(id)).Err()
}

// @func: createKey
// @date: 2023-12-04 22:37:41
// @brief: 生成第一页列表缓存在Redis中的key
// @author: Kewin Li
// @receiver r
// @param userId
// @return string
func (r *RedisArticleCache) createFirstPageKey(userId int64) string {
	return fmt.Sprintf("article:first_page:%d", userId)
}

// @func: createPreCacheKey
// @date: 2023-12-06 00:40:35
// @brief: 生成详情预加载缓存在Redis中的key
// @author: Kewin Li
// @receiver r
// @param artId
// @return string
func (r *RedisArticleCache) createPreCacheKey(artId int64) string {
	return fmt.Sprintf("article:precache:%d", artId)
}
