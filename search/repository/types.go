package repository

import (
	"context"
	"kitbook/search/domain"
)

type UserRepository interface {
	InputUser(ctx context.Context, user domain.User) error
	SearchUser(ctx context.Context, keywords []string) ([]domain.User, error)
}

type ArticleRepository interface {
	InputArticle(ctx context.Context, art domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}
