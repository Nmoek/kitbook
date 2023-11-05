// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	v := ioc.InitGinMiddlewares()
	db := ioc.InitDB()
	userDao := dao.NewGormUserDao(db)
	cmdable := InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewCacheUserRepository(userDao, userCache)
	userService := service.NewNormalUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewcodeRepository(codeCache)
	smsService := ioc.InitSmsService()
	codeService := service.NewPhoneCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebService(v, userHandler)
	return engine
}
