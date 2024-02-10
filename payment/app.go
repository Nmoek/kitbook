package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"kitbook/pkg/grpcx"
)

type App struct {
	rpcServer *grpcx.Server
	webServer *gin.Engine
	corn      *cron.Cron
}
