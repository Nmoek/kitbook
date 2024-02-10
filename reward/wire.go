//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/reward/events"
	"kitbook/reward/grpc"
	"kitbook/reward/ioc"
	"kitbook/reward/repository"
	"kitbook/reward/repository/cache"
	"kitbook/reward/repository/dao"
	"kitbook/reward/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitConsumers,
	ioc.InitSaramaClient,
	ioc.InitEtcd,

	// payment远程调用客户端
	ioc.InitPaymentClient,
	// account远程调用客户端
	ioc.InitAccountClient,
)

var rewardSvcSet = wire.NewSet(
	dao.NewGormRewardDao,
	cache.NewRedisRewardCache,
	repository.NewWechatNativeRewardRepository,
	service.NewWechatNativeRewardService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		rewardSvcSet,

		events.NewPaymentEventConsumer,
		grpc.NewRewardServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
