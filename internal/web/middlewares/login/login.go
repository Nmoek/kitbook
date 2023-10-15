// Package login
// @Description: 登录功能middleware
package login

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MiddlewareBuilder struct {
}

// @func: Build
// @date: 2023-10-09 03:10:02
// @brief: build模式-登录校验
// @author: Kewin Li
// @receiver builder
// @return gin.HandlerFunc
func (builder *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//注册、登录操作不能进行登录校验
		path := ctx.Request.URL.Path
		if path == "/users/login" || path == "/users/signup" {
			//TODO: 校验是否注册
			return
		}

		session := sessions.Default(ctx)
		if session.Get("userID") == nil {
			// 中断
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}

	}
}
