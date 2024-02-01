package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"kitbook/internal/events"
	"kitbook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	rpcServer *grpcx.Server
	webServer *gin.Engine
	corn      *cron.Cron
}
