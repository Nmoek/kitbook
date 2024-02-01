//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/payment/grpc"
	"kitbook/payment/repository"
	"kitbook/payment/repository/cache"
	"kitbook/payment/repository/dao"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
	//InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var paymentSvcSet = wire.NewSet(
	dao.NewGormPaymentDao,
	cache.NewRedisPaymentCache,
	repository.NewNativePaymentRepository,

	//wechat.NewNativePaymentService,
)

func NewPaymentServiceServer() *grpc.NewPaymentServiceServer {
	wire.Build(
		thirdPartySet,
		paymentSvcSet,
		grpc.NewPaymentServiceServer,
	)
	return new(grpc.NewPaymentServiceServer)
}

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		paymentSvcSet,

		events.NewPaymentReadEventConsumer,
		grpc.NewPaymentServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
