// Package login
// @Description: 登录功能middleware
package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "kitbook/internal/web/jwt"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	jwtHdl ijwt.JWTHandler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.JWTHandler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		jwtHdl: jwtHdl,
	}
}

// @func: Build
// @date: 2023-10-09 03:10:02
// @brief: 登录校验-普通JWT版
// @author: Kewin Li
// @receiver builder
// @return gin.HandlerFunc
func (builder *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//注册、登录操作不能进行登录校验
		if CheckIsSignupOrLogin(ctx) {
			// TODO: 日志埋点, 打印当前直接返回的URL
			return
		}

		// 校验短token
		tokenStr := builder.jwtHdl.ExtractToken(ctx)

		var claims ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			// 可以根据一些条件动态计算key, 练习方便起见使用了固定key
			return []byte(ijwt.TokenPrivateKey), nil
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
		if claims.ExpiresAt.Sub(time.Now()) < 3*time.Minute {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenStr, err = newToken.SignedString([]byte(ijwt.TokenPrivateKey))
			if err != nil {
				fmt.Printf("token SignedString err! %s \n", err)
			}

			ctx.Header("x-jwt-token", tokenStr)

		}
		// 将jwt的payload缓存在上下文中
		ctx.Set("user_token", claims)
	}
}

// @func: CheckLogin_LongShortTOken
// @date: 2023-11-14 22:52:52
// @brief: 登录校验-长短token版
// @author: Kewin Li
// @receiver builder
// @return gin.HandlerFunc
func (builder *LoginJWTMiddlewareBuilder) CheckLogin_LongShortToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//注册、登录操作不能进行登录校验
		if CheckIsSignupOrLogin(ctx) {
			// TODO: 日志埋点, 打印当前直接返回的URL
			return
		}

		// 校验短token
		tokenStr := builder.jwtHdl.ExtractToken(ctx)

		var claims ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			// 可以根据一些条件动态计算key, 练习方便起见使用了固定key
			return []byte(ijwt.TokenPrivateKey), nil
		})

		if err != nil {

			// token解析不出来
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		// token解析出来了，但是是伪造的或者过期的
		if token == nil || !token.Valid {
			// TODO: 短token过期需要生成新的短token

			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}

		userAgent := ctx.GetHeader("User-Agent")
		if userAgent == "" || userAgent != claims.UserAgent {
			//TODO: 1. 后续讲到监控警告设计时，这里还有"埋点"操作。2. 这里可以替换为检测浏览器指纹
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = builder.jwtHdl.CheckSsid(ctx, claims.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 宽松判断, 无用标识已存在则说明已经登出, redis崩溃的情形依旧可用
		//if cnt > 0 {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		// ！！！注意: 使用长短token后不再去刷新短token的过期时间
		//token是否即将过期, 是则要刷新
		//通过token中的exp time 去倒计时判断
		//if claims.ExpiresAt.Sub(time.Now()) < 3*time.Minute {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * time.Minute))
		//	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		//
		//	tokenStr, err = newToken.SignedString([]byte(web.TokenPrivateKey))
		//	if err != nil {
		//		fmt.Printf("token SignedString err! %s \n", err)
		//	}
		//
		//	ctx.Header("x-jwt-token", tokenStr)
		//
		//}

		// 将jwt的payload缓存在上下文中
		ctx.Set("user_token", claims)
	}
}
