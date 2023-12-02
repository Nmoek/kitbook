package startup

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"time"
)

func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//TODO: 打开监视器会引发test框架panic
	//monitor := &event.CommandMonitor{
	//	Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
	//		fmt.Printf("[mongo command msg] %v \n", startedEvent.Command)
	//	},
	//}

	//opts := options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor)
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client.Database("kitbook")
}
