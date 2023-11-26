// Package repository
// @Description: 帖子转发层-制作库
package repository

import (
	"context"
	"kitbook/internal/domain"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}
