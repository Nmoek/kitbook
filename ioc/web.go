// Package ioc
// @Description: Web服务组装
package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kitbook/internal/web"
	"kitbook/internal/web/middlewares"
	"strings"
	"time"
)

func InitWebService(middlewares []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {

	server := gin.Default()
	server.Use(middlewares...)
	userHdl.UserRegisterRoutes(server)
	return server
}

func InitGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true, //是否允许cookie
			AllowHeaders:     []string{"Content-Type", "authorization"},
			ExposeHeaders:    []string{"x-jwt-token"}, //允许外部访问后端的头部字段
			//AllowOrigins:     []string{"http://localhost:3000"},  //单独枚举指定
			AllowOriginFunc: func(origin string) bool {
				// 允许本机调试
				if strings.Contains(origin, "localhost") {
					return true
				}

				return strings.Contains(origin, "xxx.com.cn") //只允许公司域名
			},
			MaxAge: 12 * time.Hour,
		}),
		(&middlewares.LoginJWTMiddlewareBuilder{}).CheckLogin(),
		//TODO: 限流器中间件

	}

}
