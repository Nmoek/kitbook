package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/tag/domain"
	"strconv"
	"time"
)

type TagCache interface {
	Append(ctx context.Context, uid int64, tags ...domain.Tag) error
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
}

type RedisTagCache struct {
	client redis.Cmdable

	expiration time.Duration
}

func NewRedisTagCache(client redis.Cmdable) TagCache {
	return &RedisTagCache{
		client:     client,
		expiration: time.Minute * 10,
	}
}

func (r *RedisTagCache) Append(ctx context.Context, uid int64, tags ...domain.Tag) error {
	key := r.createKey(uid)
	pip := r.client.Pipeline()
	for _, tag := range tags {
		val, err := json.Marshal(tag)
		if err != nil {
			return err
		}

		pip.HMSet(ctx, key, strconv.FormatInt(tag.Id, 10), val).Err()
	}
	pip.Expire(ctx, key, r.expiration)
	_, err := pip.Exec(ctx)
	return err
}

func (r *RedisTagCache) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := r.client.HGetAll(ctx, r.createKey(uid)).Result()
	if err != nil {
		return nil, err
	}
	if len(res) <= 0 {
		return nil, ErrKeyNotExist
	}

	tags := make([]domain.Tag, 0, len(res))
	for _, val := range res {
		var tag domain.Tag
		err = json.Unmarshal([]byte(val), &tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil

}

func (r *RedisTagCache) createKey(uid int64) string {
	return fmt.Sprintf("tag:user_tag:%d", uid)
}
