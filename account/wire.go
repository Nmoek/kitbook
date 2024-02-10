//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/account/grpc"
	"kitbook/account/ioc"
	"kitbook/account/repository"
	"kitbook/account/repository/dao"
	"kitbook/account/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
)

var accountSvcSet = wire.NewSet(
	dao.NewGormAccountDao,
	repository.NewCacheAccountRepository,
	service.NewaccountService,
	//wechat.NewNativePaymentService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		accountSvcSet,

		grpc.NewAccountServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
