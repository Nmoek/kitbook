package domain

// 帖子状态
const (
	// 未知
	ArticleStatusUnknow = iota
	// 未发表
	ArticleStatusUnpublished
	// 已发表
	ArticleStatusPublished
	// 仅自己可见
	ArticleStatusPrivate
)

type Article struct {
	Id      int64
	Title   string
	Content string
	Status  int8
}
