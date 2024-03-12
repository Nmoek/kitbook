package dao

import "context"

type UserDao interface {
	InputUser(ctx context.Context, user User) error
	SearchUser(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleDao interface {
	InputArticle(ctx context.Context, art Article) error
	SearchArticle(ctx context.Context, artIds []int64, keywords []string) ([]Article, error)
}

type TagDao interface {
	SearchTag(ctx context.Context, uid int64, biz string, keywords []string) ([]int64, error)
}

type AnyDao interface {
	Input(ctx context.Context, index string, docID string, data string) error
}
