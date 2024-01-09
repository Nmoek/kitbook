//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/interactive/grpc"
	"kitbook/interactive/repository"
	"kitbook/interactive/repository/cache"
	"kitbook/interactive/repository/dao"
	"kitbook/interactive/service"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
	//InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGORMInteractiveDao,
	cache.NewRedisInteractiveCache,
	repository.NewArticleInteractiveRepository,
	service.NewArticleInteractiveService,
)

func NewInteractiveServer() *grpc.InteractiveServiceServer {
	wire.Build(
		thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
	)
	return new(grpc.InteractiveServiceServer)
}
