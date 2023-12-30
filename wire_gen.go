// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/internal/events/article"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/internal/web/jwt"
	"kitbook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	limiter := ioc.InitLimiter(cmdable)
	jwtHandler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, limiter, jwtHandler, logger)
	db := ioc.InitDB(logger)
	userDao := dao.NewGormUserDao(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDao, userCache)
	userService := service.NewNormalUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewcodeRepository(codeCache)
	smsService := ioc.InitSmsService(limiter)
	codeService := service.NewPhoneCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, jwtHandler, logger)
	wechatService := ioc.InitWechatService()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, jwtHandler, logger)
	articleDao := dao.NewGormArticleDao(db)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewCacheArticleRepository(articleDao, articleCache, userRepository)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewNormalArticleService(articleRepository, producer, logger)
	interactiveDao := dao.NewGORMInteractiveDao(db)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewArticleInteractiveRepository(interactiveDao, interactiveCache, logger)
	interactiveService := service.NewArticleInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(articleService, interactiveService, logger)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(interactiveRepository, client, logger)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	rankingService := service.NewBatchRankingService(interactiveService, articleService)
	rlockClient := ioc.InitRlockClient(cmdable)
	rankingJob := ioc.InitRankingJob(rankingService, rlockClient, logger)
	cron := ioc.InitJobs(logger, rankingJob)
	app := &App{
		server:    engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

// wire.go:

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDao, cache.NewRedisInteractiveCache, repository.NewArticleInteractiveRepository, service.NewArticleInteractiveService)

var rankingSvcSet = wire.NewSet(cache.NewRedisRankingCache, repository.NewCacheRankingRepository, service.NewBatchRankingService)
