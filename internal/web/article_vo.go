// Package web
// @Description: 为不直接暴露domain数据设立中间层
package web

import (
	"kitbook/internal/domain"
	"time"
)

type ArticleVo struct {
	Id         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Content    string `json:"content,omitempty"`
	AuthorId   int64  `json:"authorId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`
}

func ConvertArticleVo(art *domain.Article, isAbstract bool) ArticleVo {

	var txt string
	if isAbstract {
		txt = art.CreateAbstract()
	} else {
		txt = art.Content
	}

	vo := ArticleVo{
		Id:         art.Id,
		Title:      art.Title,
		Abstract:   txt,
		Content:    txt, // 列表展示没有必要全部内容返回
		AuthorId:   art.Author.Id,
		AuthorName: "",
		Status:     art.Status.ToUint8(),
		Ctime:      art.Ctime.Format(time.DateTime),
		Utime:      art.Utime.Format(time.DateTime),
	}
	return vo
}

func ConvertArticleVos(arts []domain.Article, isAbstract bool) []ArticleVo {
	artsVo := make([]ArticleVo, len(arts))
	for i, art := range arts {
		artsVo[i] = ConvertArticleVo(&art, isAbstract)
	}

	return artsVo

}
