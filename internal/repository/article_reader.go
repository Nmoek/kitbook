// Package repository
// @Description: 帖子转发层-线上库
package repository

import (
	"context"
	"kitbook/internal/domain"
)

type ArticleReaderRepository interface {
	// insert or update 语义
	Save(ctx context.Context, art domain.Article) error
}
