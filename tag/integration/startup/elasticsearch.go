package startup

import (
	olivere "github.com/olivere/elastic/v7"
)

func InitES() *olivere.Client {

	client, err := olivere.NewClient(olivere.SetURL("localhost:9200"),
		olivere.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return client
}
