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
	ijwt "kitbook/internal/web/jwt"
	"kitbook/ioc"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 第三方依赖
		thirdPartySet,

		dao.NewGormUserDao,
		dao.NewGormArticleDao,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		//cache.NewLocalCodeCache,

		repository.NewCacheUserRepository,
		repository.NewcodeRepository,
		repository.NewCacheArticleRepository,

		//  TODO: 如何使用多个不同的限流器
		ioc.InitLimiter,
		ioc.InitSmsService,
		service.NewNormalUserService,
		service.NewPhoneCodeService,
		InitWechatService, //不需要真的开启
		service.NewNormalArticleService,

		ioc.InitGinMiddlewares,
		ijwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitWebService,
	)

	return gin.Default()
}

func NewArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,

		dao.NewGormArticleDao,
		repository.NewCacheArticleRepository,
		service.NewNormalArticleService,
		web.NewArticleHandler,
	)

	return &web.ArticleHandler{}
}
