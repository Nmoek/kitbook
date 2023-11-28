package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
}

type ArticleStatus uint8

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

// 帖子状态
const (
	// 未知
	ArticleStatusUnknow = iota
	// 未发表
	ArticleStatusUnpblished
	// 已发表
	ArticleStatusPblished
	// 仅自己可见
	ArticleStatusPrivate
)

type Author struct {
	Id   int64
	Name string
}
