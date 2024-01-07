//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	repository2 "kitbook/interactive/repository"
	cache2 "kitbook/interactive/repository/cache"
	dao2 "kitbook/interactive/repository/dao"
	service2 "kitbook/interactive/service"
	"kitbook/internal/events/article"
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
	InitSaramaClient,
	InitSyncProducer,
	InitConsumers,
)

var interactiveSvcSet = wire.NewSet(
	dao2.NewGORMInteractiveDao,
	cache2.NewRedisInteractiveCache,
	repository2.NewArticleInteractiveRepository,
	service2.NewArticleInteractiveService,
)

var userSvcProvider = wire.NewSet(
	dao.NewGormUserDao,
	cache.NewRedisUserCache,
	repository.NewCacheUserRepository,
	service.NewNormalUserService,
)

func InitWebServer() *gin.Engine {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		interactiveSvcSet,

		dao.NewGormUserDao,
		dao.NewGormArticleDao,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,
		//cache.NewLocalCodeCache,

		repository.NewCacheUserRepository,
		repository.NewcodeRepository,
		repository.NewCacheArticleRepository,

		article.NewSaramaSyncProducer,

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
		ioc.InitWebServer,
	)

	return gin.Default()
}

func NewArticleHandler(dao dao.ArticleDao) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		interactiveSvcSet,

		article.NewSaramaSyncProducer,

		cache.NewRedisArticleCache,
		repository.NewCacheArticleRepository,
		service.NewNormalArticleService,
		web.NewArticleHandler,
	)

	return &web.ArticleHandler{}
}

func NewInteractiveService() service2.InteractiveService {
	wire.Build(
		thirdPartySet,
		interactiveSvcSet,
	)
	return service2.NewArticleInteractiveService(nil, nil)
}
