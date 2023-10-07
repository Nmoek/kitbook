package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	"net/http"
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
	group.POST("/login", h.Login)    //处理登录
	group.POST("/signup", h.SignUp)  //处理注册
}

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
		ctx.String(http.StatusInternalServerError, "系统错误！")
		return
	}

	if !isVail {
		ctx.String(http.StatusBadRequest, "邮箱格式错误！[xxx@qq.com]")
		return

	}

	// TODO: 查询数据库。密码作比对

	ctx.String(http.StatusOK, "登录成功!")

}

func (h *UserHandler) SignUp(ctx *gin.Context) {

	//1. 从Json--->结构体 协议解析
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		fmt.Printf("signup param parse fail! \n")
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

	//TODO: 写注册信息到数据库
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统错误！")
		return
	}

	ctx.String(http.StatusOK, "注册成功! ")

}

func (h *UserHandler) Profile(ctx *gin.Context) {

	type ProfileReq struct {
	}

	ctx.String(http.StatusOK, "查看用户信息!")
}
