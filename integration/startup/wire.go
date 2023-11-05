//go:build wireinject

package startup

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
		InitRedis,
		//ioc.InitFreeCache,

		dao.NewGormUserDao,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		//cache.NewLocalCodeCache,

		repository.NewCacheUserRepository,
		repository.NewcodeRepository,

		ioc.InitSmsService,
		service.NewNormalUserService,
		service.NewPhoneCodeService,

		ioc.InitGinMiddlewares,
		web.NewUserHandler,
		ioc.InitWebService,
	)

	return gin.Default()
}
