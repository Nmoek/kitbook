//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/reward/grpc"
	"kitbook/reward/repository"
	"kitbook/reward/repository/cache"
	"kitbook/reward/repository/dao"
	"kitbook/reward/service"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
	InitEtcd,
	InitPaymentClient,
	InitAccountClient,
	//InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var rewardSvcSet = wire.NewSet(
	dao.NewGormRewardDao,
	cache.NewRedisRewardCache,
	repository.NewWechatNativeRewardRepository,
	service.NewWechatNativeRewardService,
)

func NewRewardServiceServer() *grpc.RewardServiceServer {
	wire.Build(
		thirdPartySet,
		rewardSvcSet,
		grpc.NewRewardServiceServer,
	)
	return new(grpc.RewardServiceServer)
}
