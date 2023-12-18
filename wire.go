//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/internal/events/article"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/ioc"
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGORMInteractiveDao,
	cache.NewRedisInteractiveCache,
	repository.NewArticleInteractiveRepository,
	service.NewArticleInteractiveService,
)

func InitWebServer() *App {

	wire.Build(
		// 第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,
		//ioc.InitFreeCache,

		interactiveSvcSet,

		article.NewSaramaSyncProducer,
		article.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,

		dao.NewGormUserDao,
		dao.NewGormArticleDao,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,
		//cache.NewLocalCodeCache,

		repository.NewCacheUserRepository,
		repository.NewcodeRepository,
		repository.NewCacheArticleRepository,

		//  TODO: 如何使用多个不同的限流器
		ioc.InitLimiter,
		ioc.InitSmsService,
		ioc.InitWechatService,
		service.NewNormalUserService,
		service.NewPhoneCodeService,
		service.NewNormalArticleService,

		ioc.InitGinMiddlewares,
		ijwt.NewRedisJWTHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
