package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	rss "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	rdb "github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	smsTencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/config"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"kitbook/internal/service/sms"
	"kitbook/internal/service/sms/local"
	"kitbook/internal/service/sms/tencent"
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
	//smsService := initSmsTencent()
	smsService := local.NewService()

	codeSvc := initCodeService(cmd, smsService)

	initUserHandler(db, cmd, codeSvc, server)

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

// TODO: 完善腾讯云短信服务初始化
func initSmsTencent() sms.Service {

	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对secretId，secretKey。
	 * 这里采用的是从环境变量读取的方式，需要在环境变量中先设置这两个值。
	 * 你也可以直接在代码中写死密钥对，但是小心不要将代码复制、上传或者分享给他人，
	 * 以免泄露密钥对危及你的财产安全。
	 * SecretId、SecretKey 查询: https://console.cloud.tencent.com/cam/capi*/
	credential := common.NewCredential(
		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"xxx",
		"xxx",
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()
	/* SDK默认使用POST方法。
	 * 如果你一定要使用GET方法，可以在这里设置。GET方法无法处理一些较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"
	/* SDK有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	// cpf.HttpProfile.ReqTimeout = 5
	/* SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务
	 * 则必须手动指定域名，例如sms的上海金融区域名： sms.ap-shanghai-fsi.tencentcloudapi.com */
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	/* SDK默认用TC3-HMAC-SHA256进行签名，非必要请不要修改这个字段 */
	cpf.SignMethod = "HmacSHA1"
	/* 实例化要请求产品(以sms为例)的client对象
	 * 第二个参数是地域信息，可以直接填写字符串ap-guangzhou，支持的地域列表参考https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8 */
	client, _ := smsTencent.NewClient(credential, "ap-guangzhou", cpf)

	return tencent.NewService(client, "666666", "signature")
}

func initUserHandler(db *gorm.DB, cmd rdb.Cmdable, codeSvc *service.CodeService, server *gin.Engine) {
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmd)
	repo := repository.NewUserRepository(userDao, userCache)
	svc := service.NewUserService(repo)
	user := web.NewUserHandler(svc, codeSvc)
	user.UserRegisterRoutes(server)
}

func initCodeService(cmd rdb.Cmdable, sms sms.Service) *service.CodeService {
	// TODO: 生成腾讯云短信服务句柄
	codeCache := cache.NewCodeCache(cmd)
	codeRepo := repository.NewCodeRepository(codeCache)
	return service.NewCodeService(codeRepo, sms)

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
