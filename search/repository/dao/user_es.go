package dao

import (
	"context"
	"encoding/json"
	olivere "github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

type ElasticSearchUserDao struct {
	client *olivere.Client
}

func NewElasticSearchUserDao(client *olivere.Client) UserDao {
	return &ElasticSearchUserDao{
		client: client,
	}
}

func (e *ElasticSearchUserDao) InputUser(ctx context.Context, user User) error {
	_, err := e.client.Index().Index(UserIndexName).
		Id(strconv.FormatInt(user.Id, 10)). // 类似主键
		BodyJson(user).Do(ctx)
	return err
}

func (e *ElasticSearchUserDao) SearchUser(ctx context.Context, keywords []string) ([]User, error) {
	queryString := strings.Join(keywords, " ")

	query := olivere.NewMatchQuery("nikename", queryString)

	resp, err := e.client.Search(UserIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]User, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var user User
		err = json.Unmarshal(hit.Source, &user)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}

	return res, nil
}

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id,omitempty"`
	NikeName string `json:"nike_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
