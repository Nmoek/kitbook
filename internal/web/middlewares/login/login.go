// Package login
// @Description: 登录功能middleware
package login

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

		// !! 坑点：Gin中session键值对设定是覆盖式的，需要重新赋值一遍不相关字段
		session := sessions.Default(ctx)
		userID := session.Get("userID")
		if userID == nil {
			// 中断
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}

		// session超时时间刷新
		nowTime := time.Now()
		val := session.Get("update_time")
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || nowTime.Sub(lastUpdateTime) > time.Minute {

			// 对该结构体注册进行序列化
			gob.Register(time.Now())
			// TODO: 键值对赋值需要优化
			session.Set("userID", userID)

			//坑点: Gin中没有对redis的键值对设置进行字节序列化
			session.Set("update_time", nowTime)

			if err := session.Save(); err != nil {
				fmt.Printf("gin session kv param save err! %s \n", err)
			}
		}

	}
}
