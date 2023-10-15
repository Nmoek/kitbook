package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/internal/repository"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/web"
	"kitbook/internal/web/middlewares/login"
	"strings"
	"time"
)

func main() {

	db := initDB()
	server := initWebServer()

	initUserHandler(db, server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	d := dao.NewUserDao(db)
	repo := repository.NewUserRepository(d)
	svc := service.NewUserService(repo)
	user := web.NewUserHandler(svc)
	user.UserRegisterRoutes(server)
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	//初始化seesion
	loginMiddleware := login.MiddlewareBuilder{}
	store := cookie.NewStore([]byte("secret"))

	// TODO: seesionID直接放入了cookie, 这样不安全但简单起见先这么处理
	//加入登录校验middleware
	// 注意区分: 连接层sessionID 与 业务层userID
	server.Use(sessions.Sessions("userID", store), loginMiddleware.CheckLogin())

	return server
}
