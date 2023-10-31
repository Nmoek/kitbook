//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/ioc"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 第三方依赖
		ioc.InitDB,
		ioc.InitRedis,

		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		ioc.InitSmsService,
		service.NewUserService,
		service.NewCodeService,

		ioc.InitGinMiddlewares,
		web.NewUserHandler,
		ioc.InitWebService,
	)

	return gin.Default()
}
