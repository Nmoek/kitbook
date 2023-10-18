// Package login
// @Description: 登录功能middleware
package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"kitbook/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

// @func: Build
// @date: 2023-10-09 03:10:02
// @brief: build模式-登录校验 JWT版
// @author: Kewin Li
// @receiver builder
// @return gin.HandlerFunc
func (builder *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//注册、登录操作不能进行登录校验
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			//TODO: 校验是否注册
			return
		}

		// Bearer token
		authData := ctx.GetHeader("Authorization")
		if authData == "" {
			// Authorization字段未设置
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authData, " ")
		if len(segs) < 2 {
			// token格式有误
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}

		tokenStr := segs[1]

		var claims web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			// 可以根据一些条件动态计算key, 练习方便起见使用了固定key
			return []byte(web.TokenPrivateKey), nil
		})

		if err != nil {

			// token解析不出来
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		// token解析出来了，但是是伪造的或者过期的
		if token == nil || !token.Valid {

			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}

		userAgent := ctx.GetHeader("User-Agent")
		if userAgent == "" || userAgent != claims.UserAgent {
			//TODO: 1. 后续讲到监控警告设计时，这里还有"埋点"操作。2. 这里可以替换为检测浏览器指纹
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//token是否即将过期, 是则要刷新
		//通过token中的exp time 去倒计时判断
		if claims.ExpiresAt.Sub(time.Now()) < 1*time.Minute {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenStr, err = newToken.SignedString([]byte(web.TokenPrivateKey))
			if err != nil {
				fmt.Printf("token SignedString err! %s \n", err)
			}

			ctx.Header("x-jwt-token", tokenStr)

		}
		// 将jwt的payload缓存在上下文中
		ctx.Set("user_token", claims)
	}
}
