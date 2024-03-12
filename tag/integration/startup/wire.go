//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/tag/grpc"
	"kitbook/tag/repository"
	"kitbook/tag/repository/cache"
	"kitbook/tag/repository/dao"
	"kitbook/tag/service"
)

var thirdPartySet = wire.NewSet(
	InitLogger,
	InitES,
	InitDB,
	InitRedis,
	//InitConsumers,
)

var tagSvcSet = wire.NewSet(
	dao.NewGormTagDao,
	cache.NewRedisTagCache,
	repository.NewCacheTagRepository,
	service.NewTagService,
)

func NewTagServiceServer() *grpc.TagServiceServer {
	wire.Build(
		thirdPartySet,
		tagSvcSet,
		grpc.NewTagServiceServer,
	)
	return new(grpc.TagServiceServer)
}
