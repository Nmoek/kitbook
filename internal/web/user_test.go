// Package web
// @Description: 用户模块-单元测试
package web

import (
	"bytes"
	"encoding/json"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	svcmocks "kitbook/internal/service/mocks"
	"kitbook/internal/web/jwt"
	jwtmocks "kitbook/internal/web/jwt/mocks"
	"kitbook/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// UserHandlerSuite
// @Description: user测试套件
type UserHandlerSuite struct {
	suite.Suite
}

// @func: TestUserHandler_SignUp
// @date: 2024-01-11 00:56:20
// @brief: 单元测试-web接口-注册
// @author: Kewin Li
// @param t
func (u *UserHandlerSuite) TestSignUp() {
	t := u.T()
	testCases := []struct {
		// 测试条目
		name string

		// mock服务
		mock func(ctrl *gomock.Controller) service.UserService

		// 构造请求, 预期中的输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期中的输出
		wantCode int
		wantBody string
	}{
		// 注册成功
		{
			name: "Registration Success",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Ljk741610",
				}).Return(nil)

				return userSvc
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
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
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
			wantBody: "系统错误",
		},
		// 邮箱格式错误
		{
			name: "Email format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
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
		// 两次输入密码不一致
		{
			name: "Inconsistency between two passwords",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
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
		// 密码格式错误
		{
			name: "Password format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				return userSvc
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
		// 同一用户重复注册
		{
			name: "Repeated registrations by the same user",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Ljk741610",
				}).Return(service.ErrDuplicateUser)

				return userSvc
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
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)

				userSvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					Return(errors.New("其他未知错误"))

				return userSvc
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
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造handler
			userSvc := tc.mock(ctrl)
			h := NewUserHandler(userSvc, nil, nil, logger.NewNopLogger())

			// 准备服务器, 注册路由
			server := gin.Default()
			h.RegisterRoutes(server)

			// 准备请求
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 开启Web服务 并执行请求
			server.ServeHTTP(recorder, req)

			// 断言判断结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())

		})

	}
}

// @func: TestUserHandler_SignUp
// @date: 2024-01-11 00:57:45
// @brief: 单元测试-web接口-查看个人信息
// @author: Kewin Li
// @param t
func (u *UserHandlerSuite) TestProfile() {
	t := u.T()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		// 请求中入参不会变化
		//requestBuilder func(t *testing.T) *http.Request

		wantCode int
		wantErr  error
		wantBody string
		wantRes  ProfileVo
	}{
		// 查看个人信息成功
		{
			name: "Check Personal Information Successfully",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Profile(gomock.Any(), gomock.Any()).
					Return(domain.User{
						Id:       1,
						Email:    "1@qq.com",
						Phone:    "123456",
						Password: "hik12345",
						Nickname: "kewin",
						Birthday: time.Now(),
						AboutMe:  "im 666",
						Ctime:    time.Now(),
					}, nil)

				return svc

			},

			wantCode: http.StatusOK,
			wantRes: ConvertsProfileVo(&domain.User{
				Id:       1,
				Email:    "1@qq.com",
				Phone:    "123456",
				Password: "hik12345",
				Nickname: "kewin",
				Birthday: time.Now(),
				AboutMe:  "im 666",
				Ctime:    time.Now(),
			}),
		},
		// 用户不存在时进行查询
		{
			name: "Query when user does not exist",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Profile(gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrInvalidUserAccess)

				return svc

			},

			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Profile(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("其他未知错误"))

				return svc

			},

			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tc.mock(ctrl)
			h := NewUserHandler(svc, nil, nil, logger.NewNopLogger())

			// 创建服务器
			server := gin.Default()

			// 放入用户认证
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user_token", jwt.UserClaims{
					UserID: 123,
				})
			})
			h.RegisterRoutes(server)

			// 构造请求
			req, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行请求
			server.ServeHTTP(recorder, req)

			// 解析结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if len(tc.wantBody) <= 0 {
				var res ProfileVo
				err = json.NewDecoder(recorder.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRes, res)
			} else {
				assert.Equal(t, tc.wantBody, recorder.Body.String())
			}

		})
	}

}

// @func: TestEdit
// @date: 2024-01-11 02:18:58
// @brief: 单元测试-web接口-修改个人信息
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestEdit() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		requestBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody string
	}{
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/edit",
					bytes.NewReader([]byte(`{
"nickname": "name",
"birthday": "2024-01-11",
"aboutMe": "tes`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusBadRequest,
			wantBody: "系统错误",
		},
		// 生日格式错误
		{
			name: "Birthday format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/edit",
					bytes.NewReader([]byte(`{
"nickname": "name",
"birthday": "202asdsdas",
"aboutMe": "test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantBody: "非法生日格式。例: 2023-10-11",
		},
		// 修改个人信息成功
		{
			name: "Personal Information Modification Successful",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(nil)

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/edit",
					bytes.NewReader([]byte(`{
"nickname": "name",
"birthday": "2024-01-11",
"aboutMe": "test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantBody: "修改个人信息成功!",
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Edit(gomock.Any(), gomock.Any()).Return(errors.New("其他未知错误"))

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/edit",
					bytes.NewReader([]byte(`{
"nickname": "name",
"birthday": "2024-01-11",
"aboutMe": "test"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tc.mock(ctrl)
			h := NewUserHandler(svc, nil, nil, logger.NewNopLogger())

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user_token", jwt.UserClaims{
					UserID: 123,
				})
			})

			h.RegisterRoutes(server)

			req := tc.requestBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())

		})
	}

}

// @func: TestLogin
// @date: 2024-01-11 03:17:49
// @brief: 单元测试-web接口-通过session方式登录
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestLogin() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		requestBuilder func(t *testing.T) *http.Request

		after func(t *testing.T)

		wantCode int
		wantRes  Result
	}{
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/test/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "2024-01-`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			after: func(t *testing.T) {},

			wantCode: http.StatusBadRequest,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 邮箱格式错误
		{
			name: "Email format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)
				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/test/login",
					bytes.NewReader([]byte(`{
"email": "1@123123",
"password": "password"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "邮箱格式错误！[xxx@qq.com]",
			},
		},
		// 用户名或密码错误
		{
			name: "Incorrect username or password",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrInvalidUserOrPassword)

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/test/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "用户名或密码错误!",
			},
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) service.UserService {
				svc := svcmocks.NewMockUserService(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("其他未知错误"))

				return svc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/test/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// TODO: 登录成功
		//		{
		//			name: "Login successful, token set failed",
		//			mock: func(ctrl *gomock.Controller) service.UserService {
		//				svc := svcmocks.NewMockUserService(ctrl)
		//
		//				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
		//					Return(domain.User{
		//						Id: 123,
		//					}, nil)
		//
		//				return svc
		//			},
		//			requestBuilder: func(t *testing.T) *http.Request {
		//				req, err := http.NewRequest(http.MethodPost,
		//					"/users/test/login",
		//					bytes.NewReader([]byte(`{
		//"email": "1@qq.com",
		//"password": "Ljk741610"
		//}`)))
		//
		//				req.Header.Set("Content-Type", "application/json")
		//				assert.NoError(t, err)
		//
		//				return req
		//			},
		//			wantCode: http.StatusOK,
		//			wantRes: Result{
		//				Msg: "登录成功!",
		//			},
		//		},
		// 用户名或密码不正确
		// TODO: 登录成功,session信息设置失败
		//		{
		//			name: "Login successful, token set failed",
		//			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
		//				svc := svcmocks.NewMockUserService(ctrl)
		//				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
		//
		//				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
		//					Return(domain.User{}, nil)
		//
		//				ijwt.EXPECT().SetTokenWithSsid(gomock.Any(), gomock.Any()).
		//					Return(errors.New("token设置失败"))
		//
		//				return svc, ijwt
		//			},
		//			requestBuilder: func(t *testing.T) *http.Request {
		//				req, err := http.NewRequest(http.MethodPost, "/users/login",
		//					bytes.NewReader([]byte(`{
		//"email": "1@qq.com",
		//"password": "Ljk741610"
		//}`)))
		//
		//				req.Header.Set("Content-Type", "application/json")
		//				assert.NoError(t, err)
		//
		//				return req
		//			},
		//			wantCode: http.StatusOK,
		//			wantRes: Result{
		//				Msg: "系统错误",
		//			},
		//		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc := tc.mock(ctrl)
			h := NewUserHandler(userSvc, nil, nil, logger.NewNopLogger())
			server := gin.Default()
			h.RegisterRoutes(server)

			req := tc.requestBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err := json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

// @func: TestLoginWithJWT
// @date: 2024-01-11 03:17:49
// @brief: 单元测试-web接口-通过JWT方式登录
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestLoginWithJWT() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler)

		requestBuilder func(t *testing.T) *http.Request

		wantCode int
		wantRes  Result
	}{
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "2024-01-`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 邮箱格式错误
		{
			name: "Email format error",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@123123",
"password": "password"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "邮箱格式错误, 例[xxx@qq.com]",
			},
		},
		// 登录成功,token设置失败
		{
			name: "Login successful, token set failed",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, nil)

				ijwt.EXPECT().SetTokenWithSsid(gomock.Any(), gomock.Any()).
					Return(errors.New("token设置失败"))

				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 登录成功
		{
			name: "Login successful, token set failed",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, nil)

				ijwt.EXPECT().SetTokenWithSsid(gomock.Any(), gomock.Any()).
					Return(nil)

				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "登录成功!",
			},
		},
		// 用户名或密码不正确
		{
			name: "Incorrect username or password",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrInvalidUserOrPassword)

				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "用户名或密码错误!",
			},
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) (service.UserService, jwt.JWTHandler) {
				svc := svcmocks.NewMockUserService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				svc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("其他未知错误"))

				return svc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/login",
					bytes.NewReader([]byte(`{
"email": "1@qq.com",
"password": "Ljk741610"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, ijwtHdl := tc.mock(ctrl)
			h := NewUserHandler(userSvc, nil, ijwtHdl, logger.NewNopLogger())
			server := gin.Default()
			h.RegisterRoutes(server)

			req := tc.requestBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err := json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

// @func: SendLoginSMSCode
// @date: 2024-01-11 23:26:10
// @brief: 单元测试-web接口-发送手机验证码
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestSendLoginSMSCode() {
	t := u.T()
	const reqUrl = "/users/login_sms/code/send"

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		requestBuilder func(t *testing.T) *http.Request

		wantCode int
		wantRes  Result
	}{
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "516511
`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 手机号格式错误
		{
			name: "Phone format error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				return userSvc, codeSvc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18762561"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "手机号格式错误",
			},
		},
		// 验证码发送成功
		{
			name: "Verification Code Sent Successfully",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				return userSvc, codeSvc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18762850585"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "验证码发送成功",
			},
		},
		// 验证码发送过于频繁
		{
			name: "Verification codes are sent too frequently",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(service.ErrCodeSendTooMany)

				return userSvc, codeSvc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18762850585"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "验证码发送过于频繁，稍后再试",
			},
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)

				codeSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("其他未知错误"))

				return userSvc, codeSvc
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18762850585"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			h := NewUserHandler(userSvc, codeSvc, nil, logger.NewNopLogger())
			server := gin.Default()
			h.RegisterRoutes(server)

			req := tc.requestBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err := json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

// @func: TestLoginSMS
// @date: 2024-01-12 02:43:39
// @brief: 单元测试-web接口-通过验证码方式登录
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestLoginSMS() {
	t := u.T()
	const reqUrl = "/users/login_sms"

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService, jwt.JWTHandler)

		requestBuilder func(t *testing.T) *http.Request

		wantCode int
		wantRes  Result
	}{
		// 请求参数解析错误
		{
			name: "Request parameter parsing error",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "516511
`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 手机号格式错误
		{
			name: "Phone format error",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "12345",
"code": "6666"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "手机号格式错误",
			},
		},
		// 验证码查验过程出错
		{
			name: "Error in verification code checking process",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, errors.New("查验出错"))

				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18711110585",
"code": "4356"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 验证码不正确
		{
			name: "The verification code is not correct",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false, nil)

				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18711110585",
"code": "4356"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "验证码错误, 请重新输入",
			},
		},
		// 手机号注册/登录, 且注册/登录失败
		{
			name: "Phone Number Registration/Login, and Registration/Login Fail",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)

				userSvc.EXPECT().SignupOrLoginWithPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("注册/登录出错"))

				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18711110585",
"code": "4356"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
		// 手机号注册/登录, 且注册/登录成功
		{
			name: "Phone Number Registration/Login, and Registration/Login Successful",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)

				userSvc.EXPECT().SignupOrLoginWithPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, nil)

				ijwt.EXPECT().SetTokenWithSsid(gomock.Any(), gomock.Any()).
					Return(nil)

				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18711110585",
"code": "4356"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "登录成功",
			},
		},
		// token设置失败
		{
			name: "Failed to set the token",
			mock: func(ctrl *gomock.Controller) (
				service.UserService,
				service.CodeService,
				jwt.JWTHandler) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil)

				userSvc.EXPECT().SignupOrLoginWithPhone(gomock.Any(), gomock.Any()).
					Return(domain.User{}, nil)

				ijwt.EXPECT().SetTokenWithSsid(gomock.Any(), gomock.Any()).
					Return(errors.New("token设置出错"))

				return userSvc, codeSvc, ijwt
			},
			requestBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, reqUrl,
					bytes.NewReader([]byte(`{
"phone": "18711110585",
"code": "4356"
}`)))

				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)

				return req
			},
			wantCode: http.StatusUnauthorized,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc, ijwt := tc.mock(ctrl)
			h := NewUserHandler(userSvc, codeSvc, ijwt, logger.NewNopLogger())
			server := gin.Default()
			h.RegisterRoutes(server)

			req := tc.requestBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err := json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})
	}

}

// @func: TestRefreshToken
// @date: 2024-01-12 03:40:28
// @brief: 单元测试-web接口-长短token刷新
// @author: Kewin Li
// @receiver u
func (u *UserHandlerSuite) TestRefreshToken() {
	t := u.T()
	const reqUrl = "/users/refresh_token"

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) jwt.JWTHandler

		wantCode int
	}{
		// 获取错误token
		{
			name: "Get error token",
			mock: func(ctrl *gomock.Controller) jwt.JWTHandler {
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				ijwt.EXPECT().ExtractToken(gomock.Any()).Return("")

				return ijwt
			},
			wantCode: http.StatusUnauthorized,
		},
		// token过期或未解析正确
		{
			name: "The token has expired or is not resolved correctly",
			mock: func(ctrl *gomock.Controller) jwt.JWTHandler {
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				// 放入一个过期token
				ijwt.EXPECT().ExtractToken(gomock.Any()).
					Return("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTk5ODk4OTMsIlVzZXJJRCI6MSwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4wIiwiU3NpZCI6IjRlY2M2ZjFmLWMxOTAtNDA1YS1hNzk4LWRlNzljZjQxYjYxZCJ9.yhMtz-PO1peuDr4bsWoLufu8LtttHbSbgegaMPyQFzE")

				return ijwt
			},
			wantCode: http.StatusUnauthorized,
		},
		//其余情况并不好模拟, 以实际情况为主
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ijwt := tc.mock(ctrl)
			h := NewUserHandler(nil, nil, ijwt, logger.NewNopLogger())
			server := gin.Default()
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodGet, "/users/refresh_token", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)

		})
	}
}

// @func: TestLogout
// @date: 2024-01-12 03:59:27
// @brief: 单元测试-web接口-用户注销session方式
// @author: Kewin Li
// @receiver u
// @return func
//func (u *UserHandlerSuite) TestLogout() {
//	t := u.T()
//
//	// TODO: 如何测试session?
//	testCases := []struct {
//		name string
//
//		wantCode int
//		wantRes  Result
//	}{
//		// 用户成功登出
//		{
//			name: "User successfully logged out",
//
//			wantCode: http.StatusOK,
//			wantRes: Result{
//				Msg: "用户已退出登录",
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			h := NewUserHandler(nil, nil, nil, logger.NewNopLogger())
//			server := gin.Default()
//			h.RegisterRoutes(server)
//
//			req, err := http.NewRequest(http.MethodPost,
//				"/users/test/logout", nil)
//			assert.NoError(t, err)
//			recorder := httptest.NewRecorder()
//
//			server.ServeHTTP(recorder, req)
//
//			var res Result
//			err = json.NewDecoder(recorder.Body).Decode(&res)
//			assert.NoError(t, err)
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			assert.Equal(t, tc.wantRes, res)
//		})
//	}
//}

// @func: LogoutWithJWT
// @date: 2024-01-12 04:11:24
// @brief: 单元测试-web接口-用户注销jwt方式
// @author: Kewin Li
// @receiver u
// @return func
func (u *UserHandlerSuite) TestLogoutWithJWT() {
	t := u.T()

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) jwt.JWTHandler

		wantCode int
		wantRes  Result
	}{
		// 用户成功登出
		{
			name: "User successfully logged out",
			mock: func(ctrl *gomock.Controller) jwt.JWTHandler {
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				ijwt.EXPECT().ClearToken(gomock.Any()).Return(nil)

				return ijwt
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "用户已经退出登录",
			},
		},
		// 清理token
		{
			name: "User successfully logged out",
			mock: func(ctrl *gomock.Controller) jwt.JWTHandler {
				ijwt := jwtmocks.NewMockJWTHandler(ctrl)

				ijwt.EXPECT().ClearToken(gomock.Any()).Return(errors.New("清理出错"))

				return ijwt
			},
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ijwt := tc.mock(ctrl)
			h := NewUserHandler(nil, nil, ijwt, logger.NewNopLogger())
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user_token", jwt.UserClaims{
					UserID: 123,
				})

			})
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/logout", nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestUserHandler(t *testing.T) {
	suite.Run(t, &UserHandlerSuite{})
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

	h := NewUserHandler(nil, nil, nil, logger.NewNopLogger())

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
//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// 使用mock 生成的初始化服务
//	userSvc := svcmocks.NewMockUserService(ctrl)
//
//	// 设置模拟场景
//	userSvc.EXPECT().Signup(gomock.Any(), domain.User{
//		Id:    1,
//		Email: "123@qq.com",
//	}).Return(errors.New("这是一个mock测试"))
//
//	err := userSvc.Signup(context.Background(), domain.User{
//		Id:    1,
//		Email: "123@qq.com",
//	})
//
//	t.Log(err)
//}
