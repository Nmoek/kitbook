package dao

import "context"

type UserDao interface {
	InputUser(ctx context.Context, user User) error
	SearchUser(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleDao interface {
	InputArticle(ctx context.Context, art Article) error
	SearchArticle(ctx context.Context, keywords []string) ([]Article, error)
}
