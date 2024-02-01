package main

import (
	"github.com/gin-gonic/gin"
	"kitbook/internal/events"
	"kitbook/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	rpcServer *grpcx.Server
	webServer *gin.Engine
}
