package dao

import "context"

type ArticleAuthorDao interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}
