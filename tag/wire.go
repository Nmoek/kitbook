//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/tag/grpc"
	"kitbook/tag/ioc"
	"kitbook/tag/repository"
	"kitbook/tag/repository/cache"
	"kitbook/tag/repository/dao"
	"kitbook/tag/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitES,
	ioc.InitDB,
	ioc.InitRedis,
)

var tagSvcSet = wire.NewSet(
	dao.NewGormTagDao,
	cache.NewRedisTagCache,
	repository.NewCacheTagRepository,
	service.NewTagService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		tagSvcSet,

		grpc.NewTagServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
