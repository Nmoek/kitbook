// Package web
// @Description: 用户模块-单元测试
package web

import (
	"bytes"
	"context"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	svcmocks "kitbook/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		// 测试条目
		name string

		// mock服务
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 构造请求, 预期中的输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期中的输出
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Ljk741610",
				}).Return(nil)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "Ljk741610",
"confirmPassword": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusOK,
			wantBody: "注册成功！",
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "Ljk741610
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusBadRequest,
			wantBody: "参数解析错误！",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@",
"password": "Ljk741610",
"confirmPassword": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误！[xxx@qq.com]",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "Ljk741610",
"confirmPassword": "Ljk"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusOK,
			wantBody: "两次密码输入不一致！",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "741610",
"confirmPassword": "741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusOK,
			wantBody: "必须包含大小写字母和数字的组合，不能使用特殊字符，长度在8-16之间",
		},
		{
			name: "重复注册",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)

				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Ljk741610",
				}).Return(service.ErrDuplicateUser)

				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/users/signup",
					bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "Ljk741610",
"confirmPassword": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},

			wantCode: http.StatusOK,
			wantBody: service.ErrDuplicateUser.Error(),
		},
	}

	for _, tc := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// 构造handler
		userSvc, codeSvc := tc.mock(ctrl)
		h := NewUserHandler(userSvc, codeSvc)

		// 准备服务器, 注册路由
		server := gin.Default()
		h.UserRegisterRoutes(server)

		// 准备请求
		req := tc.reqBuilder(t)
		recorder := httptest.NewRecorder()

		// 开启Web服务 并执行请求
		server.ServeHTTP(recorder, req)

		// 断言判断结果
		assert.Equal(t, tc.wantCode, recorder.Code)
		assert.Equal(t, tc.wantBody, recorder.Body.String())
	}
}

// @func: TestHttp
// @date: 2023-11-04 04:00:04
// @brief: http接口测试
// @author: Kewin Li
// @param t
func TestHttp(t *testing.T) {

	// TODO: 难点是如何构造一个Body的方法
	_, err := http.NewRequest(http.MethodPost,
		"/users/signup",
		bytes.NewReader([]byte("this is body")))
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	assert.Equal(t, t, http.StatusOK, recorder.Code)
}

// @func: TestEmailRegExp
// @date: 2023-11-04 04:09:04
// @brief: 测试邮箱校验接口
// @author: Kewin Li
// @param t
func TestEmailRegExp(t *testing.T) {

	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "不带@",
			email: "854981891",
			match: false,
		},
		{
			name:  "不带后缀.com",
			email: "24231@",
			match: false,
		},
		{
			name:  "存在重复@",
			email: "123@163@qq.com",
			match: false,
		},
		{
			name:  "合法邮箱",
			email: "123@qq.com",
			match: true,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, val := range testCases {
		tc := val
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRegExp.MatchString(tc.email)
			assert.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})

	}

}

// @func: TestPhoneRegExp
// @date: 2023-10-29 21:33:08
// @brief: 测试手机号校验
// @author: Kewin Li
// @param t
func TestPhoneRegExp(t *testing.T) {

	const phoneRegexPattern = "(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}"
	phoneRegExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	ok, err := phoneRegExp.MatchString("15662850585")

	t.Logf("ok:%v, err:%v \n", ok, err)
}

// @func: TestMock
// @date: 2023-11-04 17:09:38
// @brief: mock使用入门测试
// @author: Kewin Li
// @param t
func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 使用mock 生成的初始化服务
	userSvc := svcmocks.NewMockUserService(ctrl)

	// 设置模拟场景
	userSvc.EXPECT().Signup(gomock.Any(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	}).Return(errors.New("这是一个mock测试"))

	err := userSvc.Signup(context.Background(), domain.User{
		Id:    1,
		Email: "123@qq.com",
	})

	t.Log(err)
}
