package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"kitbook/ioc"
	"net/http"
	"time"
)

func main() {
	// 初始化配置模块
	initViper()
	initPrometheus()
	tpCancel := ioc.InitOTEL()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// 控制关闭超时
		tpCancel(ctx)

	}()

	// 初始化Web服务
	app := InitApp()
	// 开始热榜定时任务
	app.cron.Start()
	// 等待热榜定时任务退出
	defer func() {
		<-app.cron.Stop().Done()
	}()

	server := app.server
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	//server := gin.Default()
	//
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello, this is K8s!!")
	//	return
	//})

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}

}

func initViperV1() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initViper() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/dev.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initPrometheus() {
	go func() {
		// 专门给prometheus用的端口
		http.Handle("metrics", promhttp.Handler())
		err := http.ListenAndServe("localhost:8082", nil)
		if err != nil {
			panic(err)
		}
	}()
}

//func initUserHandler(db *gorm.DB, cmd rdb.Cmdable, codeSvc *service.PhoneCodeService, server *gin.Engine) {
//	userDao := dao.NewGormUserDao(db)
//	userCache := cache.NewRedisUserCache(cmd)
//	repo := repository.NewCacheUserRepository(userDao, userCache)
//	svc := service.NewNormalUserService(repo)
//	user := web.NewUserHandler(svc, codeSvc)
//	user.RegisterRoutes(server)
//}

//func userSession(server *gin.Engine) {
//
//	//初始化seesion
//	loginMiddleware := middlewares.LoginMidwodlewareBuilder{}
//
//	// 1. 使用cookie存储session
//	//store := cookie.NewStore([]byte("secret"))
//
//	// 2. 使用memstore存储session; 第一个密钥用于身份认证, 第二个密钥用于数据加解密
//	//store := memstore.NewStore([]byte("tHaegpgS1uxjmH3E9suduGmXECFm7CEk"), []byte("s6AjedURwVItfEsrhKS4QKvAUnRWJCcL"))
//
//	// 3. 使用redis存储session
//	store, err := rss.NewStore(10, "tcp", config.Config.Redis.Addr, "",
//		[]byte("tHaegpgS1uxjmH3E9suduGmXECFm7CEk"),
//		[]byte("s6AjedURwVItfEsrhKS4QKvAUnRWJCcL"))
//
//	// 4. 其他的store介质
//
//	if err != nil {
//		fmt.Printf("redis store err! %s \n", err)
//		os.Exit(-1)
//	}
//
//	// TODO: seesionID直接放入了cookie, 这样不安全但简单起见先这么处理
//	//加入登录校验middleware
//	// 注意区分: 连接层sessionID 与 业务层userID
//	server.Use(sessions.Sessions("sessionID", store), loginMiddleware.CheckLogin())
//
//}

//func userJWT(server *gin.Engine) {
//	jwtMiddlewareBuilder := middlewares.LoginJWTMiddlewareBuilder{}
//
//	server.Use(jwtMiddlewareBuilder.CheckLogin())
//}
