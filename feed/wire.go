//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/feed/grpc"
	"kitbook/feed/ioc"
	"kitbook/feed/repository"
	"kitbook/feed/repository/dao"
	"kitbook/feed/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcd,
	ioc.InitFollowClient,
	ioc.InitHandler,
)

var feedSvcSet = wire.NewSet(
	dao.NewFeedPushEventDao,
	dao.NewFeedPullEventDao,
	repository.NewFeedEventRepository,
	service.NewFeedService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		feedSvcSet,

		grpc.NewFeedServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
