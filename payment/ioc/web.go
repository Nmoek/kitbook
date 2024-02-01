package ioc

import (
	"github.com/gin-gonic/gin"
	"kitbook/payment/web"
)

func InitWebServer(wechatHdl *web.WeChatNativeHandler) *gin.Engine {
	server := gin.Default()
	wechatHdl.RegisterRoutes(server)
	return server
}
