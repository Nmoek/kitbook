// Package login
// @Description: 登录校验通用组件
package login

import "github.com/gin-gonic/gin"

// 所有注册、登录的URL
var signupOrLoginPaths = []string{
	"/users/login",
	"/users/signup",
	"/users/login_sms",
	"/users/login_sms/code/send",
}

// @func: checkIsSignupOrLogin
// @date: 2023-10-30 22:18:30
// @brief: 判断当前操作是否属于注册、登录之一
// @author: Kewin Li
// @param ctx
// @return bool
func CheckIsSignupOrLogin(ctx *gin.Context) bool {
	requestPath := ctx.Request.URL.Path

	for _, path := range signupOrLoginPaths {
		if requestPath == path {
			return true
		}
	}

	return false
}
