package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	rss "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	rdb "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/config"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/internal/web/middlewares/login"
	"os"
	"strings"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()
	cmd := initRedis()

	initUserHandler(db, cmd, server)

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

func initUserHandler(db *gorm.DB, cmd rdb.Cmdable, server *gin.Engine) {
	d := dao.NewUserDao(db)
	c := cache.NewUserCache(cmd)
	repo := repository.NewUserRepository(d, c)
	svc := service.NewUserService(repo)
	user := web.NewUserHandler(svc)
	user.UserRegisterRoutes(server)
}

func initDB() *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// !middleware注册
	server.Use(cors.New(cors.Config{
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
	}))

	//  TODO: redis使用的限流器

	//userSession(server)
	userJWT(server)

	return server
}

func initRedis() rdb.Cmdable {

	return rdb.NewClient(&rdb.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password, // no password docs
		DB:       0,                            // use default DB
	})
}

func userSession(server *gin.Engine) {

	//初始化seesion
	loginMiddleware := login.LoginMiddlewareBuilder{}

	// 1. 使用cookie存储session
	//store := cookie.NewStore([]byte("secret"))

	// 2. 使用memstore存储session; 第一个密钥用于身份认证, 第二个密钥用于数据加解密
	//store := memstore.NewStore([]byte("tHaegpgS1uxjmH3E9suduGmXECFm7CEk"), []byte("s6AjedURwVItfEsrhKS4QKvAUnRWJCcL"))

	// 3. 使用redis存储session
	store, err := rss.NewStore(10, "tcp", config.Config.Redis.Addr, "",
		[]byte("tHaegpgS1uxjmH3E9suduGmXECFm7CEk"),
		[]byte("s6AjedURwVItfEsrhKS4QKvAUnRWJCcL"))

	// 4. 其他的store介质

	if err != nil {
		fmt.Printf("redis store err! %s \n", err)
		os.Exit(-1)
	}

	// TODO: seesionID直接放入了cookie, 这样不安全但简单起见先这么处理
	//加入登录校验middleware
	// 注意区分: 连接层sessionID 与 业务层userID
	server.Use(sessions.Sessions("sessionID", store), loginMiddleware.CheckLogin())

}

func userJWT(server *gin.Engine) {
	jwtMiddlewareBuilder := login.LoginJWTMiddlewareBuilder{}

	server.Use(jwtMiddlewareBuilder.CheckLogin())
}
