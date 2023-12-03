package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

func ToArticleStatus(status uint8) ArticleStatus {
	return ArticleStatus(status)
}

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
