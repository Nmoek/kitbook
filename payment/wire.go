//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/payment/events"
	"kitbook/payment/grpc"
	"kitbook/payment/ioc"
	"kitbook/payment/repository"
	"kitbook/payment/repository/cache"
	"kitbook/payment/repository/dao"
	"kitbook/payment/web"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitSaramaClient,

	ioc.InitWechatClient,
	ioc.InitWechatConfig,
	ioc.InitWechatNotifyHandler,

	ioc.InitJobs,
	ioc.InitSyncWechatOrderJob,
	ioc.InitRlockClient,
	ioc.InitSyncProducer,
	//InitConsumers,
)

var paymentSvcSet = wire.NewSet(
	dao.NewGormPaymentDao,
	cache.NewRedisPaymentCache,
	repository.NewNativePaymentRepository,
	ioc.InitWechatNativeService,
	//wechat.NewNativePaymentService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		paymentSvcSet,

		events.NewSaramaSyncProducer,
		grpc.NewPaymentServiceServer,
		ioc.InitGRpcServer,
		web.NewWeChatNativeHandler,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}