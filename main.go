package main

func main() {

	server := InitWebServer()

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

//func initUserHandler(db *gorm.DB, cmd rdb.Cmdable, codeSvc *service.PhoneCodeService, server *gin.Engine) {
//	userDao := dao.NewGormUserDao(db)
//	userCache := cache.NewRedisUserCache(cmd)
//	repo := repository.NewCacheUserRepository(userDao, userCache)
//	svc := service.NewNormalUserService(repo)
//	user := web.NewUserHandler(svc, codeSvc)
//	user.UserRegisterRoutes(server)
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
