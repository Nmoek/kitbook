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
	Content    string `json:"content,omitempty"`
	AuthorId   int64  `json:"authorId,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     uint8  `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`
}

func ConvertArticleVo(art *domain.Article) ArticleVo {
	return ArticleVo{
		Id:         art.Id,
		Title:      art.Title,
		Content:    art.Content,
		AuthorId:   art.Author.Id,
		AuthorName: "",
		Status:     art.Status.ToUint8(),
		Ctime:      art.Ctime.Format(time.DateTime),
		Utime:      art.Utime.Format(time.DateTime),
	}
}

func ConvertArticleVos(arts []domain.Article) []ArticleVo {
	artsVo := make([]ArticleVo, len(arts))
	for i, art := range arts {
		artsVo[i] = ConvertArticleVo(&art)
	}

	return artsVo

}
