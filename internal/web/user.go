package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	"net/http"
	"time"
)

// 正则表达式
const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	//长度大于8位小于16位，大小写密码组合，不包含特殊字符的密码校验
	passwordRegexPattern = "^(?=.*\\d)(?=.*[a-z])(?=.*[A-Z]).{8,16}$"
	// 手机号校验
	phoneRegexPattern = "(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}"
)

const TokenPrivateKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C5"

const bizLogin = "login"

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	phoneRegExp    *regexp.Regexp
	svc            service.UserService
	code           service.CodeService
}

// @func: NewUserHandler
// @date: 2023-10-06 18:36:02
// @brief: 创建用户模块句柄
// @author: Kewin Li
// @return *UserHandler
func NewUserHandler(svc service.UserService, code service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRegExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		svc:            svc,
		code:           code,
	}
}

// @func: UserRegisterRoutes
// @date: 2023-10-04 20:36:47
// @brief: 用户模块-路由注册
// @author: Kewin Li
// @receiver h
// @param server
func (h *UserHandler) UserRegisterRoutes(server *gin.Engine) {
	// 注册、登录等功能对应的URL(路由规则)
	group := server.Group("/users")
	group.GET("/profile", h.Profile)     //查询用户信息
	group.POST("/login", h.LoginWithJWT) //JWT登录
	//group.POST("/login", h.Login)    //登录
	group.POST("/signup", h.SignUp) //注册
	group.POST("/edit", h.Edit)     //修改个人信息

	//手机验证码登录相关对应URL
	group.POST("login_sms/code/send", h.SendLoginSMSCode)
	group.POST("login_sms", h.LoginSMS)
}

// @func: setSession
// @date: 2023-10-29 23:09:23
// @brief: 设置Session
// @author: Kewin Li
// @receiver h
// @param ctx
// @param id
// @return error
func (h *UserHandler) setSession(ctx *gin.Context, id int64) error {
	session := sessions.Default(ctx)
	session.Set("userID", id)
	session.Options(sessions.Options{
		HttpOnly: true,
		Secure:   true,
		MaxAge:   900, //15min
	})

	return session.Save()
}

// @func: Login
// @date: 2023-10-12 03:30:44
// @brief: 用户模块-登录
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *UserHandler) Login(ctx *gin.Context) {

	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 邮箱、密码检测。防止进入到数据库中检索比对，拖慢系统
	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		fmt.Printf("login param parse fail! \n")
		return
	}

	//2. 文本校验--正则表达式
	isVail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}

	if !isVail {
		ctx.String(http.StatusOK, "邮箱格式错误！[xxx@qq.com]")
		return
	}

	user, err := h.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:

		// 设置Seesion
		err = h.setSession(ctx, user.Id)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		ctx.String(http.StatusOK, "登录成功！")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不正确!")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

}

// @func: setJWTToken
// @date: 2023-10-29 23:06:05
// @brief: 设置JWT
// @author: Kewin Li
// @receiver h
// @param ctx
// @param id
func (h *UserHandler) setJWTToken(ctx *gin.Context, id int64) {
	// 设置JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:    id,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	})

	tokenStr, err := token.SignedString([]byte(TokenPrivateKey))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误!")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)
}

// @func: LoginWithJWT
// @date: 2023-10-16 18:55:54
// @brief: 用户模块-登录通过JWT方式
// @author: Kewin Li
// @receiver h
// @param context
func (h *UserHandler) LoginWithJWT(ctx *gin.Context) {

	type UserReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req UserReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	isValid, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isValid {
		ctx.String(http.StatusOK, "邮箱格式错误, 例[xxx@qq.com]")
		return
	}

	user, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:

		// 设置JWT
		h.setJWTToken(ctx, user.Id)

		ctx.String(http.StatusOK, "登录成功!")

	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码错误!")
	default:
		ctx.String(http.StatusOK, "系统错误!")
	}

}

// @func: SignUp
// @date: 2023-10-12 03:31:10
// @brief: 用户模块-注册
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *UserHandler) SignUp(ctx *gin.Context) {

	//1. 从Json--->结构体 协议解析
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "参数解析错误！")
		return
	}

	//2. 文本校验--正则表达式
	isVail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "系统错误！")
		ctx.String(http.StatusOK, "系统错误！")
		return
	}

	if !isVail {
		//ctx.String(http.StatusBadRequest, "邮箱格式错误！[xxx@qq.com]")
		ctx.String(http.StatusOK, "邮箱格式错误！[xxx@qq.com]")
		return

	}

	// 两次密码不一致检测
	if req.Password != req.ConfirmPassword {
		//ctx.String(http.StatusBadRequest, "两次密码输入不一致！")
		ctx.String(http.StatusOK, "两次密码输入不一致！")
		return
	}

	isVail, err = h.passwordRegExp.MatchString(req.Password)
	if err != nil {
		fmt.Printf("密码校验出错！\n")
		//ctx.String(http.StatusInternalServerError, "系统错误！")
		ctx.String(http.StatusOK, "系统错误！")

		return
	}

	if !isVail {
		//ctx.String(http.StatusBadRequest, "必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-16之间")
		ctx.String(http.StatusOK, "必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-16之间")
		return
	}

	//写注册信息到数据库
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功！")
	case service.ErrDuplicateUser:
		ctx.String(http.StatusOK, "%s", service.ErrDuplicateUser)
	default:
		ctx.String(http.StatusOK, "系统错误!")
	}

}

func checkBySession(ctx *gin.Context) int64 {
	// 通过sessionID拿到是哪一个用户
	session := sessions.Default(ctx)
	if session.Get("userID") == nil {
		return -1
	}
	return session.Get("userID").(int64)
}

func checkByJWT(ctx *gin.Context) int64 {
	val := ctx.MustGet("user_token")
	if val == nil {
		return -1
	}
	claims := val.(UserClaims)
	return claims.UserID

}

// @func: Edit
// @date: 2023-10-12 03:31:51
// @brief: 用户模块-修改个人信息
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *UserHandler) Edit(ctx *gin.Context) {

	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "参数解析错误！")
		return
	}

	//1. 数据校验
	//2. 如何查询当前修改信息的用户是谁?
	//userID := checkBySession(ctx)
	userID := checkByJWT(ctx)
	if userID < 0 {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "非法生日格式。例: 2023-10-11")
		return
	}

	err = h.svc.Edit(ctx, domain.User{
		Id:       userID,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})

	if err != nil {
		fmt.Printf("[Edit] %s \n", err)
		ctx.String(http.StatusOK, "系统错误!")
		return
	}

	ctx.String(http.StatusOK, "修改个人信息成功!")
}

// @func: Profile
// @date: 2023-10-14 17:50:47
// @brief: 用户模块-查看个人信息
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *UserHandler) Profile(ctx *gin.Context) {

	type UserResponse struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}

	// 1. 用户ID
	//userID := checkBySession(ctx)
	userID := checkByJWT(ctx)
	if userID < 0 {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	user, err := h.svc.Profile(ctx, userID)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, UserResponse{
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    "",
			Birthday: user.Birthday.Format(time.DateOnly),
			AboutMe:  user.AboutMe,
		})
	case service.ErrInvalidUserAccess:
		ctx.String(http.StatusOK, "非法用户访问!")
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}
}

// @func: SendLoginSMSCode
// @date: 2023-10-29 04:15:20
// @brief: 用户模块-向用户发送手机验证码
// @author: Kewin Li
// @receiver h
// @param context
func (h *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type SendLoginSMSCodeReq struct {
		Phone string `json:"phone"`
	}

	var req SendLoginSMSCodeReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 手机号校验
	isValid, err := h.phoneRegExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isValid {
		ctx.String(http.StatusOK, "手机号格式错误")
		return
	}

	// TODO: 使用一个本地调试，不需要真的使用短信服务
	err = h.code.Send(ctx, bizLogin, req.Phone)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码发送过于频繁，稍后再试",
		})

	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}

}

// @func: LoginSMS
// @date: 2023-10-29 04:15:47
// @brief: 用户模块-用户验证码登录
// @author: Kewin Li
// @receiver h
// @param context
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type LoginSMSReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req LoginSMSReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 手机号校验
	isValid, err := h.phoneRegExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isValid {
		ctx.String(http.StatusOK, "手机号格式错误")
		return
	}

	ok, err := h.code.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		// TODO: 验证码登录错误, 日志埋点
		//fmt.Printf("code login fail! %s \n", err)
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统异常",
		})

		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码错误, 请重新输入",
		})

		return
	}

	//发现当前手机号没有进行注册，需要提示注册并进入注册流程
	// TODO: 如果此时的手机号已经注册过邮箱，如何将该手机号和邮箱合并？合并时应注意什么问题

	user, err := h.svc.SignupOrLoginWithPhone(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		return
	}

	// 设置JWT
	h.setJWTToken(ctx, user.Id)

	// 设置Seesion 该实现也可以
	//err = h.setSession(ctx, user.Id)
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})

}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}
