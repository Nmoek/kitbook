package main

import (
	"github.com/gin-gonic/gin"
	"kitbook/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
