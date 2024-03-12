package dao

import (
	"context"
	olivere "github.com/olivere/elastic/v7"
)

type anyDao struct {
	client *olivere.Client
}

func NewAnyDao(client *olivere.Client) AnyDao {
	return &anyDao{
		client: client,
	}
}

func (a *anyDao) Input(ctx context.Context, index string, docID string, data string) error {
	_, err := a.client.Index().Id(index).Id(docID).BodyString(data).Do(ctx)
	return err
}
