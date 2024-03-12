package dao

import (
	"context"
	"encoding/json"
	olivere "github.com/olivere/elastic/v7"
)

const TagsIndexName = "tags_index"

type ElasticsearchTagDao struct {
	client olivere.Client
}

func (e *ElasticsearchTagDao) SearchTag(ctx context.Context, uid int64, biz string, keywords []string) ([]int64, error) {
	query := olivere.NewBoolQuery().Must(
		olivere.NewTermQuery("uid", uid),
		olivere.NewTermQuery("biz", biz),
		olivere.NewTermsQueryFromStrings("tags", keywords...),
	)

	// 查询出所有命中标签的资源ID
	resp, err := e.client.Search(TagsIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var bizTag BizTags
		err = json.Unmarshal(hit.Source, &bizTag)
		if err != nil {
			return nil, err
		}
		res = append(res, bizTag.BizId)
	}

	return res, nil

}

type BizTags struct {
	Uid   int64    `json:"uid,omitempty"`
	Biz   string   `json:"biz,omitempty"`
	BizId int64    `json:"biz_id,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}
