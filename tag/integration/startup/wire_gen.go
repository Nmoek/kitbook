// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/tag/grpc"
	"kitbook/tag/repository"
	"kitbook/tag/repository/cache"
	"kitbook/tag/repository/dao"
	"kitbook/tag/service"
)

// Injectors from wire.go:

func NewTagServiceServer() *grpc.TagServiceServer {
	db := InitDB()
	tagDao := dao.NewGormTagDao(db)
	cmdable := InitRedis()
	tagCache := cache.NewRedisTagCache(cmdable)
	tagRepository := repository.NewCacheTagRepository(tagDao, tagCache)
	tagService := service.NewTagService(tagRepository)
	tagServiceServer := grpc.NewTagServiceServer(tagService)
	return tagServiceServer
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitLogger,
	InitES,
	InitDB,
	InitRedis,
)

var tagSvcSet = wire.NewSet(dao.NewGormTagDao, cache.NewRedisTagCache, repository.NewCacheTagRepository, service.NewTagService)
