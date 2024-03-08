package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/follow/domain"
	"strconv"
)

const (
	// 粉丝数
	fieldFollowerCnt = "follower_cnt"
	// 关注数
	fieldFolloweeCnt = "followee_cnt"
)

type FollowCache interface {
	Follow(ctx context.Context, follower int64, followee int64) error
	CancelFollow(ctx context.Context, follower int64, followee int64) error

	SetStaticsInfo(ctx context.Context, uid int64, s domain.FollowStatics) error
	GetStaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

type RedisFollowCache struct {
	client redis.Cmdable
}

func NewRedisFollowCache(client redis.Cmdable) FollowCache {
	return &RedisFollowCache{
		client: client,
	}
}

// @func: Follow
// @date: 2024-02-12 22:25:23
// @brief: 关注某个用户-关注数、粉丝数同步增加
// @author: Kewin Li
// @receiver r
// @param ctx
// @param follower
// @param followee
// @return error
func (r *RedisFollowCache) Follow(ctx context.Context, follower int64, followee int64) error {
	return r.updateStaticsInfo(ctx, follower, followee, 1)
}

// @func: CancelFollow
// @date: 2024-02-12 22:25:56
// @brief: 取消关注某个用户-关注数、粉丝数同步减少
// @author: Kewin Li
// @receiver r
// @param ctx
// @param follower
// @param followee
// @return error
func (r *RedisFollowCache) CancelFollow(ctx context.Context, follower int64, followee int64) error {
	return r.updateStaticsInfo(ctx, follower, followee, -1)

}

// @func: updateStaticsInfo
// @date: 2024-02-12 22:26:05
// @brief: 更改关注数、粉丝数真正实现
// @author: Kewin Li
// @receiver r
// @param ctx
// @param follower
// @param followee
// @param delta
// @return error
func (r *RedisFollowCache) updateStaticsInfo(ctx context.Context, follower int64, followee int64, delta int64) error {
	tx := r.client.TxPipeline()
	// 使用hash表

	// 增加被关注者的粉丝数
	r.client.HIncrBy(ctx, r.createKey(follower), fieldFolloweeCnt, delta)
	// 增加关注者的关注数量
	r.client.HIncrBy(ctx, r.createKey(followee), fieldFollowerCnt, delta)

	_, err := tx.Exec(ctx)
	return err
}

// @func: SetStaticsInfo
// @date: 2024-02-12 22:30:07
// @brief: 设置某个用户的关注数+粉丝数
// @author: Kewin Li
// @receiver r
// @param ctx
// @param uid
// @param followerCnt
// @param followeeCnt
// @return error
func (r *RedisFollowCache) SetStaticsInfo(ctx context.Context, uid int64, s domain.FollowStatics) error {
	return r.client.HMSet(ctx, r.createKey(uid), fieldFollowerCnt, s.Followers, fieldFolloweeCnt, s.Followees).Err()
}

// @func: StaticsInfo
// @date: 2024-02-12 22:25:08
// @brief: 获取某个用户的关注数+粉丝数
// @author: Kewin Li
// @receiver r
// @param ctx
// @param uid
// @return domain.FollowStatics
// @return error
func (r *RedisFollowCache) GetStaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error) {

	res, err := r.client.HGetAll(ctx, r.createKey(uid)).Result()
	if err != nil {
		return domain.FollowStatics{}, err
	}

	if len(res) <= 0 {
		return domain.FollowStatics{}, ErrKeyNotExist
	}

	return r.convertsDomainFSFromCache(res), nil
}

func (r *RedisFollowCache) createKey(uid int64) string {
	return fmt.Sprintf("follow:statics:%d", uid)
}

// @func: convertsDomainIntrFromCache
// @date: 2023-12-15 17:52:01
// @brief: Interactive转换 redis cache ---> domain
// @author: Kewin Li
// @receiver r
// @param res
// @return domain.Interactive
func (r *RedisFollowCache) convertsDomainFSFromCache(res map[string]string) domain.FollowStatics {

	followers, _ := strconv.ParseInt(res[fieldFollowerCnt], 10, 64)
	followees, _ := strconv.ParseInt(res[fieldFolloweeCnt], 10, 64)

	return domain.FollowStatics{
		Followers: followers,
		Followees: followees,
	}
}
