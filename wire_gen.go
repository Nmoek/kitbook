// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/internal/web/jwt"
	"kitbook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
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
	engine := ioc.InitWebService(v, userHandler, oAuth2WechatHandler)
	return engine
}
