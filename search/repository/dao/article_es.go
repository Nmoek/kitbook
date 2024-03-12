package dao

import (
	"context"
	"encoding/json"
	olivere "github.com/olivere/elastic/v7"
	"kitbook/search/domain"
	"strconv"
	"strings"
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

func (e *ElasticSearchArticleDao) SearchArticle(ctx context.Context, artIds []int64, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	status := olivere.NewTermQuery("status", domain.ArticleStatusPublished) // 帖子必须可见
	title := olivere.NewMatchQuery("title", queryString)
	content := olivere.NewMatchQuery("content", queryString)

	tmp := []any{}
	for _, id := range artIds {
		tmp = append(tmp, any(id))
	}
	// 更高的权重 意味着标签命中会排在前面
	tag := olivere.NewTermsQuery("id", tmp).Boost(2)

	// Should 类似Or语义
	or := olivere.NewBoolQuery().Should(title, content, tag)
	// 等价于  `where status = ArticleStatusPublished  and (title = xx or content = xx or tag IN [] )`
	query := olivere.NewBoolQuery().Must(status, or)
	resp, err := e.client.Search(ArticleIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]Article, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var art Article
		err = json.Unmarshal(hit.Source, &art)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}

	return res, nil
}

const ArticleIndexName = "article_index"

type Article struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Status  int8   `json:"status,omitempty"`
}
