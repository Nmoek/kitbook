package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"kitbook/internal/web"
	"strings"
	"time"
)

func main() {
	//user := &web.UserHandler{}
	user := web.NewUserHandler()
	server := gin.Default()

	// TODO: 跨域问题？搁置，后续再看
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

	user.UserRegisterRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		fmt.Printf("server fun err! %s", err)
	}
}
