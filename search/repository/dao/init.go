package dao

import (
	"context"
	_ "embed"
	"fmt"
	olivere "github.com/olivere/elastic/v7"
	"golang.org/x/sync/errgroup"
	"time"
)

var (
	//go:embed user_index.json
	userIndexDef string
	//go:embed article_index.json
	artIndexDef string
	//go:embed tags_index.json
	tagsIndexDef string
)

func InitEs(client *olivere.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var eg errgroup.Group

	eg.Go(func() error {
		return tryCreateIndex(ctx, client, UserIndexName, userIndexDef)
	})

	eg.Go(func() error {

		return tryCreateIndex(ctx, client, ArticleIndexName, artIndexDef)
	})

	eg.Go(func() error {
		return tryCreateIndex(ctx, client, TagsIndexName, tagsIndexDef)
	})

	return eg.Wait()
}

func tryCreateIndex(ctx context.Context, client *olivere.Client, idxName string, idxCfg string) error {
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("%v had exists\n", idxName)
		return nil
	}

	_, err = client.CreateIndex(idxName).Body(idxCfg).Do(ctx)
	return err
}
