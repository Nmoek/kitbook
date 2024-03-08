package dao

import (
	"context"
	olivere "github.com/olivere/elastic/v7"
	"strconv"
)

type ElasticSearchArticleDao struct {
	client *olivere.Client
}

func NewElasticSearchArticleDao(client *olivere.Client) ArticleDao {
	return &ElasticSearchArticleDao{
		client: client,
	}
}

func (e *ElasticSearchArticleDao) InputArticle(ctx context.Context, art Article) error {
	_, err := e.client.Index().Index(ArticleIndexName).
		Id(strconv.FormatInt(art.Id, 10)). // 类似主键
		BodyJson(art).Do(ctx)
	return err
}

func (e *ElasticSearchArticleDao) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

const ArticleIndexName = "article_index"

type Article struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Status  int8   `json:"status,omitempty"`
}
