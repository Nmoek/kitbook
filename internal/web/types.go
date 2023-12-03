package web

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(server *gin.Engine)
}

type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
