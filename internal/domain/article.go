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

// 截取摘要的最大长度
const abstractMaxLen = 256

// @func: Abstract
// @date: 2023-12-04 22:26:53
// @brief: 通过截取长度生成摘要
// @author: Kewin Li
// @receiver a
func (a *Article) CreateAbstract() string {
	// 考虑中文
	str := []rune(a.Content)

	if len(str) > abstractMaxLen {
		str = str[:abstractMaxLen]
	}

	return string(str)
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
