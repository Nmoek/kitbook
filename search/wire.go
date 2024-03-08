//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/search/grpc"
	"kitbook/search/ioc"
	"kitbook/search/repository"
	"kitbook/search/repository/dao"
	"kitbook/search/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitLogger,
	ioc.InitES,
	//InitConsumers,
)

var syncSvcSet = wire.NewSet(
	dao.NewElasticSearchUserDao,
	dao.NewElasticSearchArticleDao,
	repository.NewArticleSyncRepository,
	repository.NewUserSyncRepository,
	service.NewSyncService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		syncSvcSet,

		grpc.NewSyncServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
