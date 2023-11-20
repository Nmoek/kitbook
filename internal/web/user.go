package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
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

const bizLogin = "login"

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	phoneRegExp    *regexp.Regexp
	svc            service.UserService
	code           service.CodeService
	jwtHdl         ijwt.JWTHandler
	l              logger.Logger
}

// @func: NewUserHandler
// @date: 2023-10-06 18:36:02
// @brief: 创建用户模块句柄
// @author: Kewin Li
// @return *UserHandler
func NewUserHandler(svc service.UserService,
	code service.CodeService,
	jwtHdl ijwt.JWTHandler,
	l logger.Logger) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		phoneRegExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		svc:            svc,
		code:           code,
		jwtHdl:         jwtHdl,
		l:              l,
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

	// 通过长token刷新短token
	group.GET("/refresh_token", h.RefreshToken)

	//手机验证码登录相关对应URL
	group.POST("login_sms/code/send", h.SendLoginSMSCode)
	group.POST("login_sms", h.LoginSMS)

	//group.POST("/logout", h.Logout)

	group.POST("/logout", h.LogoutWithJWT)

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
	var err error
	var isVail bool
	var user domain.User
	var msg = UserLogMsg{
		KeyNum: logger.LOG_LOGIN,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_LOGIN]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	//2. 文本校验--正则表达式
	isVail, err = h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		goto ERR
	}

	if !isVail {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s][%s], 邮箱格式错误", ctx.ClientIP(), req.Email)
		otherMsg := fmt.Sprintf("[%s][%s], 邮箱格式错误", ctx.ClientIP(), req.Email)
		h.l.WARN(logKey, logger.Field{"email", otherMsg})
		ctx.String(http.StatusOK, "邮箱格式错误！[xxx@qq.com]")
		return
	}

	user, err = h.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:

		// 设置Seesion
		err = h.setSession(ctx, user.Id)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			goto ERR
		}

		//msg.Level = logger.InfoLevel
		//msg.OtherMsg = fmt.Sprintf("[%s][%s], 登录成功", ctx.ClientIP(), req.Email)
		otherMsg := fmt.Sprintf("[%s][%s], 登录成功", ctx.ClientIP(), req.Email)
		h.l.INFO(logKey, logger.Field{"success", otherMsg})
		ctx.String(http.StatusOK, "登录成功！")
		return
	case service.ErrInvalidUserOrPassword:
		msg.Level = logger.WarnLevel
		ctx.String(http.StatusOK, "用户名或密码不正确!")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_LOGIN], msg)
	return

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
	var err error
	var isValid bool
	var user domain.User
	var msg = UserLogMsg{
		KeyNum: logger.LOG_LOGIN,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_LOGIN]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	isValid, err = h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	if !isValid {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 输入邮箱格式错误")
		otherMsg := fmt.Sprintf("[%s][%s], 邮箱格式错误", ctx.ClientIP(), req.Email)
		h.l.WARN(logKey, logger.Field{"email", otherMsg})
		ctx.String(http.StatusOK, "邮箱格式错误, 例[xxx@qq.com]")
		goto ERR
	}

	user, err = h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		err = h.jwtHdl.SetTokenWithSsid(ctx, user.Id)
		if err != nil {
			ctx.JSON(http.StatusOK, Result{
				Msg: "系统错误",
			})

			goto ERR
		}

		//msg.Level = logger.InfoLevel
		//msg.OtherMsg = fmt.Sprintf("[%s][%s], 登录成功", ctx.ClientIP(), req.Email)
		otherMsg := fmt.Sprintf("[%s][%s], 登录成功", ctx.ClientIP(), req.Email)
		h.l.INFO(logKey, logger.Field{"success", otherMsg})
		ctx.String(http.StatusOK, "登录成功!")
		return

	case service.ErrInvalidUserOrPassword:
		msg.Level = logger.WarnLevel
		ctx.String(http.StatusOK, "用户名或密码错误!")
	default:
		ctx.String(http.StatusOK, "系统错误!")
	}

ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_LOGIN], msg)
	return
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
	var err error = nil // 默认没有错误
	var isVail bool
	var msg = UserLogMsg{
		KeyNum: logger.LOG_SIGNUP, //具体逻辑下具体放置
		Level:  logger.ErrorLevel, //默认为错误
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_SIGNUP]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "参数解析错误！")
		goto ERR
	}

	//2. 文本校验--正则表达式
	isVail, err = h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		goto ERR
	}

	if !isVail {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 邮箱格式输入错误", ctx.ClientIP())
		ctx.String(http.StatusOK, "邮箱格式错误！[xxx@qq.com]")
		//goto ERR

		h.l.WARN(logKey, logger.Field{"email", "邮箱格式错误"})
		return
	}

	// 两次密码不一致检测
	if req.Password != req.ConfirmPassword {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 两次密码输入不一致", ctx.ClientIP())

		otherMsg := fmt.Sprintf("[%s] 两次密码输入不一致", ctx.ClientIP())
		h.l.WARN(logKey, logger.Field{"email", otherMsg})

		ctx.String(http.StatusOK, "两次密码输入不一致！")
		return
	}

	isVail, err = h.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		goto ERR
	}

	if !isVail {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 输入密码格式错误", ctx.ClientIP())

		otherMsg := fmt.Sprintf("[%s] 两次密码输入不一致", ctx.ClientIP())
		h.l.WARN(logKey, logger.Field{"email", otherMsg})

		ctx.String(http.StatusOK, "必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-16之间")
		return
		//goto ERR
	}

	//写注册信息到数据库
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:
		//msg.Level = logger.InfoLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] %s, 注册成功", ctx.ClientIP(), req.Email)
		otherMsg := fmt.Sprintf("[%s] %s, 注册成功", ctx.ClientIP(), req.Email)
		h.l.INFO(logKey, logger.Field{"success", otherMsg})
		ctx.String(http.StatusOK, "注册成功！")
		return //无报错就返回

	case service.ErrDuplicateUser:
		ctx.String(http.StatusOK, "%s", service.ErrDuplicateUser)
	default:
		ctx.String(http.StatusOK, "系统错误!")
	}

ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_SIGNUP], msg)
	return
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
	claims := val.(ijwt.UserClaims)
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
	var err error
	var birthday time.Time
	var userID int64
	var msg = UserLogMsg{
		KeyNum: logger.LOG_EDIT,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_EDIT]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "参数解析错误！")
		goto ERR
	}

	//1. 数据校验
	//2. 如何查询当前修改信息的用户是谁?
	//userID := checkBySession(ctx)
	userID = checkByJWT(ctx)
	if userID < 0 {
		err = fmt.Errorf("[%s] userID 非法", ctx.ClientIP())
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	birthday, err = time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 输入非法生日格式", ctx.ClientIP())

		otherMsg := fmt.Sprintf("[%s] 输入非法生日格式", ctx.ClientIP())
		h.l.WARN(logKey, logger.Field{"birthday", otherMsg})

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
		ctx.String(http.StatusOK, "系统错误!")
		goto ERR
	}

	h.l.INFO(logKey, logger.Field{"success", "修改个人信息成功"})
	ctx.String(http.StatusOK, "修改个人信息成功!")
	return

ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_EDIT], msg)
	return
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

	var userID int64
	var err error
	var user domain.User
	var msg = UserLogMsg{
		KeyNum: logger.LOG_PROFILE,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_PROFILE]
	// 1. 用户ID
	//userID := checkBySession(ctx)
	userID = checkByJWT(ctx)
	if userID < 0 {
		err = fmt.Errorf("[%s] userID 非法", ctx.ClientIP())
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	user, err = h.svc.Profile(ctx, userID)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, UserResponse{
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    "",
			Birthday: user.Birthday.Format(time.DateOnly),
			AboutMe:  user.AboutMe,
		})

		h.l.INFO(logKey, logger.Field{"success", "查看个人信息"})
		return
	case service.ErrInvalidUserAccess:
		err = fmt.Errorf("[%s] 非法用户访问", ctx.ClientIP())
		ctx.String(http.StatusOK, "非法用户访问!")
	default:
		ctx.String(http.StatusOK, "系统错误！")
	}

ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_PROFILE], msg)
	return
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
	var err error
	var isValid bool
	var msg = UserLogMsg{
		KeyNum: logger.LOG_SEND_SMSCODE,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_SEND_SMSCODE]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	// 手机号校验
	isValid, err = h.phoneRegExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	if !isValid {
		//msg.Level = logger.WarnLevel
		//msg.OtherMsg = fmt.Sprintf("[%s] 输入手机号格式错误")
		otherMsg := fmt.Sprintf("[%s] 输入手机号格式错误")
		h.l.WARN(logKey, logger.Field{"phone", otherMsg})
		ctx.String(http.StatusOK, "手机号格式错误")
		return
	}

	// 使用一个本地调试，不需要真的使用短信服务
	err = h.code.Send(ctx, bizLogin, req.Phone)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码发送成功",
		})
		//msg.Level = logger.InfoLevel
		//TODO: 手机号为敏感信息, 需要单独处理
		//msg.OtherMsg = fmt.Sprintf("[%s] %s, 验证码发送成功", ctx.ClientIP(), req.Phone)

		otherMsg := fmt.Sprintf("[%s] %s, 验证码发送成功", ctx.ClientIP(), req.Phone)
		h.l.INFO(logKey, logger.Field{"success", otherMsg})
		return
	case service.ErrCodeSendTooMany:
		msg.Level = logger.WarnLevel
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码发送过于频繁，稍后再试",
		})

		h.l.WARN(logKey, logger.Field{"sms_code", err.Error()})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
	}
ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_SEND_SMSCODE], msg)
	return

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
	var err error
	var isValid bool
	var ok bool
	var user domain.User
	var msg = UserLogMsg{
		KeyNum: logger.LOG_LOGINSMS,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_LOGINSMS]

	err = ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		goto ERR
	}

	// 手机号校验
	isValid, err = h.phoneRegExp.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !isValid {
		//msg.Level = logger.WarnLevel
		//TODO: 手机号为敏感信息, 需要单独处理
		//msg.OtherMsg = fmt.Sprintf("[%s] in:%s, 手机号格式错误", ctx.ClientIP(), req.Phone)

		h.l.WARN(logKey, logger.Field{"phone",
			fmt.Sprintf("[%s] in:%s, 手机号格式错误", ctx.ClientIP(), req.Phone)})
		ctx.String(http.StatusOK, "手机号格式错误")
		return
	}

	ok, err = h.code.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统异常",
		})
		goto ERR
	}

	if !ok {
		//msg.Level = logger.WarnLevel
		//TODO: 手机号为敏感信息, 需要单独处理
		//msg.OtherMsg = fmt.Sprintf("[%s][%s], 输入验证码错误", ctx.ClientIP(), req.Phone)

		h.l.WARN(logKey, logger.Field{"sms_code",
			fmt.Sprintf("[%s][%s], 输入验证码错误", ctx.ClientIP(), req.Phone)})
		ctx.JSON(http.StatusOK, Result{
			Msg: "验证码错误, 请重新输入",
		})

		return
	}

	//发现当前手机号没有进行注册，需要提示注册并进入注册流程
	// TODO: 如果此时的手机号已经注册过邮箱，如何将该手机号和邮箱合并？合并时应注意什么问题

	user, err = h.svc.SignupOrLoginWithPhone(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		goto ERR
	}

	err = h.jwtHdl.SetTokenWithSsid(ctx, user.Id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		goto ERR
	}

	// 设置Seesion 该实现也可以
	//err = h.setSession(ctx, user.Id)
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	//msg.Level = logger.InfoLevel
	//msg.OtherMsg = fmt.Sprintf("[%s][%s][%d][%s], 登录成功",
	//	ctx.ClientIP(), req.Phone, user.Id, user.Email)

	// TODO: 手机号敏感信息，单独处理
	h.l.INFO(logKey, logger.Field{"success",
		fmt.Sprintf("[%s][%s][%d][%s], 登录成功", ctx.ClientIP(), req.Phone, user.Id, user.Email)})

	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
	return
ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_LOGINSMS], msg)
	return
}

// @func: RefreshToken
// @date: 2023-11-14 00:06:57
// @brief: 获取新的短token(access_token)
// @author: Kewin Li
// @receiver h
// @param context
func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 约定将refresh_token放入在auth字段中
	tokenStr := h.jwtHdl.ExtractToken(ctx)

	var claims ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ijwt.TokenRefreshKey), nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.jwtHdl.CheckSsid(ctx, claims.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.jwtHdl.SetJWTToken(ctx, claims.UserID, claims.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "换取access_token成功",
	})

}

// @func: Logout
// @date: 2023-11-14 22:25:03
// @brief: sessions版本-用户注销(登出)
// @author: Kewin Li
// @receiver h
// @param context
func (h *UserHandler) Logout(ctx *gin.Context) {
	var msg = UserLogMsg{
		KeyNum: logger.LOG_LOGOUT,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_LOGOUT]

	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: -1, //立即删除cookie
	})

	err := session.Save()

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg: "系统错误",
		})
		goto ERR
	}

	//msg.Level = logger.InfoLevel
	//msg.OtherMsg = fmt.Sprintf("[%s] 退出登录", ctx.ClientIP())
	h.l.INFO(logKey, logger.Field{"success", fmt.Sprintf("[%s] 退出登录", ctx.ClientIP())})
	ctx.JSON(http.StatusOK, Result{
		Msg: "用户已退出登录",
	})

	return
ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_LOGOUT], msg)
	return
}

// @func: LogoutWithJWT
// @date: 2023-11-15 00:40:41
// @brief: jwt版本-用户注销(登出)
// @author: Kewin Li
// @receiver h
// @param ctx
func (h *UserHandler) LogoutWithJWT(ctx *gin.Context) {
	var msg = UserLogMsg{
		KeyNum: logger.LOG_LOGOUT,
		Level:  logger.ErrorLevel,
	}
	var logKey = logger.UserLogMsgKey[logger.LOG_LOGOUT]

	err := h.jwtHdl.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, "系统错误")
		goto ERR
	}

	//msg.Level = logger.InfoLevel
	//msg.OtherMsg = fmt.Sprintf("[%s] 退出登录", ctx.ClientIP())
	h.l.INFO(logKey, logger.Field{"success", fmt.Sprintf("[%s] 退出登录", ctx.ClientIP())})
	ctx.JSON(http.StatusOK, Result{
		Msg: "用户已经退出登录",
	})
	return
ERR:
	msg.Err = err
	ctx.Set(logger.UserLogMsgKey[logger.LOG_LOGOUT], msg)
	return
}

type UserLogMsg struct {
	KeyNum int
	// 错误应打印的基本级别
	Level    int
	Err      error
	OtherMsg string
}
