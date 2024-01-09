//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/interactive/events"
	"kitbook/interactive/grpc"
	"kitbook/interactive/ioc"
	"kitbook/interactive/repository"
	"kitbook/interactive/repository/cache"
	"kitbook/interactive/repository/dao"
	"kitbook/interactive/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitConsumers,
	ioc.InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGORMInteractiveDao,
	cache.NewRedisInteractiveCache,
	repository.NewArticleInteractiveRepository,
	service.NewArticleInteractiveService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		interactiveSvcSet,

		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
