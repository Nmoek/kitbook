package service

import (
	"context"
	"kitbook/search/domain"
)

type SyncService interface {
	InputUser(ctx context.Context, user domain.User) error
	InputArticle(ctx context.Context, art domain.Article) error
	Input(ctx context.Context, idxName string, docID string, data string) error
}

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}
