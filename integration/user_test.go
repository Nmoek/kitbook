// Package integration
// @Description: 集成测试-用户模块
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"kitbook/integration/startup"
	"kitbook/internal/web"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// @func: TestUserHandle_SnedSMSCode
// @date: 2023-11-05 22:32:14
// @brief: 集成测试-发送短信验证码
// @author: Kewin Li
// @param t
func TestUserHandle_SendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebServer()

	testCases := []struct {
		name string

		// 准备数据
		before func(t *testing.T)

		// 验证数据(中间件中的数据)
		after func(t *testing.T)

		phone string

		wantCode       int
		wantJsonBody   web.Result
		wantStringBody string
	}{

		{
			name: "验证码发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 1. redis中有数据, 且数据是否正常
				// 2. redis中数据过期时间是否正常
				// 3. 验证完需要把数据删除
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:18762850585"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				// 验证码是随机的, 并不知道验证码的具体数值, 只能知道是否存在数据
				assert.True(t, len(code) > 0)

				leftTime, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, leftTime > time.Minute*9+time.Second*50)

				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "18762850585",
			wantCode: http.StatusOK,
			wantJsonBody: web.Result{
				Msg: "验证码发送成功",
			},
		},
		{
			name: "手机号格式错误/手机号为空",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			phone:          "",
			wantCode:       http.StatusOK,
			wantStringBody: "手机号格式错误",
		},
		{
			name: "验证码发送频繁",
			before: func(t *testing.T) {
				// 1. redis中插入一个验证码, 且TTL过期时长不超过1min
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()

				key := "phone_code:login:18762850585"
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. redis中有数据, 且数据是否正常
				// 2. redis中数据过期时间是否正常
				// 3. 验证完需要把数据删除
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:18762850585"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "18762850585",
			wantCode: http.StatusOK,
			wantJsonBody: web.Result{
				Msg: "验证码发送过于频繁，稍后再试",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试前的准备数据
			tc.before(t)
			// 测试后的验证数据+清除数据
			defer tc.after(t)

			// 构造请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行, 发送请求
			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}

			if strings.Contains(recorder.Header().Get("Content-Type"), "json") {
				// 解析响应结构体
				var res web.Result
				err = json.NewDecoder(recorder.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Equal(t, tc.wantJsonBody, res)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStringBody, recorder.Body.String())
		})
	}

}
