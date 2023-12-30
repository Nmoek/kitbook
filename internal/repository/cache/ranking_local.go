package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"kitbook/internal/domain"
	"time"
)

var ErrLocalCacheInvalid = errors.New("本地缓存失效")

type LocalRankingCache struct {
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func (l *LocalRankingCache) Set(ctx context.Context, arts []domain.Article) error {

	l.topN.Store(arts)
	l.ddl.Store(time.Now().Add(l.expiration))
	return nil
}

func (l *LocalRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := l.ddl.Load()
	arts := l.topN.Load()
	if len(arts) <= 0 || ddl.Before(time.Now()) {
		return nil, ErrLocalCacheInvalid
	}

	return arts, nil
}

func (l *LocalRankingCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := l.topN.Load()
	if len(arts) <= 0 {
		return nil, ErrLocalCacheInvalid
	}

	return arts, nil
}
