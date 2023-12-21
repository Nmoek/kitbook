// Package ioc
// @Description: Web服务组装
package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"kitbook/internal/web"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/internal/web/middlewares"
	"kitbook/pkg/ginx/prometheus"
	"kitbook/pkg/limiter"
	"kitbook/pkg/logger"
	"strings"
	"time"
)

func InitWebServer(middlewares []gin.HandlerFunc,
	userHdl *web.UserHandler,
	wechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler) *gin.Engine {

	server := gin.Default()
	server.Use(middlewares...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(client redis.Cmdable,
	limiter limiter.Limiter,
	jwtHdl ijwt.JWTHandler,
	l logger.Logger) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "kewin",
		Subsystem: "kitbook",
		Name:      "gin_http",
		Help:      "统计gin中http接口数据",
	}

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
		pb.BuildResponseTIme(),
		pb.BuildActiveRequest(),
		//middlewares.NewLoginJWTMiddlewareBuilder(jwtHdl).CheckLogin(),

		//限流器中间件 1000 QPS/s
		//ratelimit.NewMiddlewareBuilder(limiter).Build(),
		// 传入要实现的日志打印
		//middlewares.NewLogMiddlewareBuilder(func(ctx context.Context, al middlewares.AccessLog) {
		//
		//	l.DEBUG("sys_log", logger.Field{Key: "req", Val: al})
		//
		//}).AllowReqBody().AllowRespBody().Build(),
		middlewares.NewLoginJWTMiddlewareBuilder(jwtHdl).CheckLogin_LongShortToken(),
	}

}
