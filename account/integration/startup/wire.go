//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/account/grpc"
	"kitbook/account/repository"
	"kitbook/account/repository/dao"
	"kitbook/account/service"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
	//InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var accountSvcSet = wire.NewSet(
	dao.NewGormAccountDao,
	repository.NewCacheAccountRepository,
	service.NewaccountService,
	//wechat.NewNativePaymentService,
)

func NewAccountServiceServer() *grpc.AccountServiceServer {
	wire.Build(
		thirdPartySet,
		accountSvcSet,
		grpc.NewAccountServiceServer,
	)
	return new(grpc.AccountServiceServer)
}
