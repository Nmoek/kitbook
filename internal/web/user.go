package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
)

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	svc            *service.UserService
}

// @func: NewUserHandler
// @date: 2023-10-06 18:36:02
// @brief: 创建用户模块句柄
// @author: Kewin Li
// @return *UserHandler
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
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
	group.GET("/profile", h.Profile) //查询用户信息
	group.POST("/login", h.Login)    //登录
	group.POST("/signup", h.SignUp)  //注册
	group.POST("/edit", h.Edit)      //修改个人信息
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
		session := sessions.Default(ctx)
		session.Set("userID", user.Id)
		session.Options(sessions.Options{
			HttpOnly: true,
			Secure:   true,
			MaxAge:   900, //15min
		})

		err = session.Save()
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
		fmt.Printf("邮箱校验出错！\n")
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
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "%s", service.ErrDuplicateEmail)
	default:
		ctx.String(http.StatusOK, "系统错误!")
	}

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

	// 通过sessionID拿到是哪一个用户
	session := sessions.Default(ctx)
	if session.Get("userID") == nil {
		// 中断
		ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}
	userID := session.Get("userID").(int64)

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
	session := sessions.Default(ctx)

	if session.Get("userID") == nil {
		ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}

	userID := session.Get("userID").(int64)

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
