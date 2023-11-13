package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"kitbook/internal/service"
	"kitbook/internal/service/oauth2/wechat"
	"net/http"
)

type OAuth2WechatHandler struct {
	wechatSvc       wechat.Service
	userSvc         service.UserService
	jwtHandler      // 通过组合方式去共享一些函数/接口
	key             string
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatSvc:       svc,
		userSvc:         userSvc,
		key:             "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C6",
		stateCookieName: "jwt-state",
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/oauth2/wechat/")
	group.GET("/authurl", o.Auth2URL)
	// 回调发送微信的code时，并不知道使用的http method
	group.Any("/callback", o.Callback)
}

// @func: Auth2URL
// @date: 2023-11-11 20:29:54
// @brief: 微信服务-返回跳转URL
// @author: Kewin Li
// @receiver o
// @param ctx
func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New().String()

	url, err := o.wechatSvc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "构造跳转URL失败",
		})
		return
	}

	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

	isValid, err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "state解析失败",
		})

		return
	}
	if !isValid {
		ctx.JSON(http.StatusOK, Result{
			Msg: "state不合法",
		})
		return
	}

	code := ctx.Query("code")

	info, err := o.wechatSvc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "微信授权失败",
		})

		return
	}

	user, err := o.userSvc.SignupOrLoginWithWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})

		return
	}

	o.setJWTToken(ctx, user.Id)

	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

// @func: setStateCookie
// @date: 2023-11-13 02:11:18
// @brief: 将state校验码设置到cookie中
// @author: Kewin Li
// @receiver o
// @param ctx
// @param state
// @return error
func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	// 设置JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
	})

	tokenStr, err := token.SignedString([]byte(o.key))
	if err != nil {
		return err
	}

	ctx.SetCookie("jwt-state", tokenStr,
		600,
		"/oauth2/wechat/callback",
		"", false, true)

	return nil
}

// @func: verifyState
// @date: 2023-11-13 02:12:28
// @brief: 校验拿到的state
// @author: Kewin Li
// @receiver o
// @param ctx
// @return bool
// @return error
func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) (bool, error) {
	state := ctx.Query("state")

	cookie, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return false, err
	}

	var claims StateClaims
	_, err = jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(o.key), nil
	})

	if err != nil {
		return false, err
	}

	if claims.State != state {
		// TODO: 日志埋点
		return false, errors.New("state不匹配")
	}

	return true, nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
