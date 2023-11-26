package dao

import "context"

type ArticleReaderDao interface {
	Upsert(crx context.Context, art Article) error
}
