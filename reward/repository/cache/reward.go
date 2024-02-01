package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"kitbook/reward/domain"
	"time"
)

type RewardCache interface {
	SetCodeUrl(ctx context.Context, codeUrl domain.CodeURL, rwd domain.Reward) error
	GetCodeUrl(ctx context.Context, rwd domain.Reward) (domain.CodeURL, error)
	DelCodeUrl(ctx context.Context, bizId int64, biz string) error
}

type RedisRewardCache struct {
	client redis.Cmdable
}

func NewRedisRewardCache(client redis.Cmdable) RewardCache {
	return &RedisRewardCache{
		client: client,
	}
}

func (r *RedisRewardCache) SetCodeUrl(ctx context.Context, codeUrl domain.CodeURL, rwd domain.Reward) error {
	val, err := json.Marshal(&codeUrl)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(rwd.Target.BizId, rwd.Target.Biz), val, time.Minute*30).Err()
}

func (r *RedisRewardCache) GetCodeUrl(ctx context.Context, rwd domain.Reward) (domain.CodeURL, error) {
	val, err := r.client.Get(ctx, r.createKey(rwd.Target.BizId, rwd.Target.Biz)).Bytes()
	if err != nil {
		return domain.CodeURL{}, err
	}

	var res domain.CodeURL
	err = json.NewDecoder(bytes.NewReader(val)).Decode(&res)
	return res, err
}

func (r *RedisRewardCache) DelCodeUrl(ctx context.Context, bizId int64, biz string) error {
	return r.client.Del(ctx, r.createKey(bizId, biz)).Err()
}

func (r *RedisRewardCache) createKey(bizId int64, biz string) string {
	return fmt.Sprintf("%s-%d", biz, bizId)
}
