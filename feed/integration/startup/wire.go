//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/search/grpc"
	"kitbook/search/repository"
	"kitbook/search/repository/dao"
	"kitbook/search/service"
)

var thirdPartySet = wire.NewSet(
	InitLogger,
	InitES,
)

var syncSvcSet = wire.NewSet(
	dao.NewElasticSearchUserDao,
	dao.NewElasticSearchArticleDao,
	repository.NewArticleSyncRepository,
	repository.NewUserSyncRepository,
	service.NewSyncService,
)

var searchSvcSet = wire.NewSet(
	dao.NewElasticSearchUserDao,
	dao.NewElasticSearchArticleDao,
	repository.NewArticleSyncRepository,
	repository.NewUserSyncRepository,
	service.NewSearchService,
)

func NewSyncServiceServer() *grpc.SyncServiceServer {
	wire.Build(
		thirdPartySet,
		syncSvcSet,
		grpc.NewSyncServiceServer,
	)
	return new(grpc.SyncServiceServer)
}

func NewSearchServiceServer() *grpc.SearchServiceServer {
	wire.Build(
		thirdPartySet,
		searchSvcSet,
		grpc.NewSearchServiceServer,
	)
	return new(grpc.SearchServiceServer)
}
