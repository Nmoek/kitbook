package dao

import (
	"context"
	olivere "github.com/olivere/elastic/v7"
	"strconv"
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
	//TODO implement me
	panic("implement me")
}

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id,omitempty"`
	NikeName string `json:"nike_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
