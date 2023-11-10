package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/service/sms"
	smsmocks "kitbook/internal/service/sms/mocks"
	"testing"
)

// @func: TestSMSService_Send
// @date: 2023-11-10 18:52:56
// @brief: 单元测试-token验证接口调用
// @author: Kewin Li
// @param t
func TestSMSService_Send(t *testing.T) {
	const testToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUcGwiOiIxMjM0NTYifQ.Dit7d-5yzhwBtFCf5sWBE9JSnsnkwFUxUlYIKtMY-3U"
	const TokenPrivateKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C5"

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (svc sms.Service, key []byte)

		wantErr error
	}{
		//
		{
			name: "token解析成功, 且发送成功",
			mock: func(ctrl *gomock.Controller) (svc sms.Service, key []byte) {
				sms := smsmocks.NewMockService(ctrl)
				sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return sms, []byte(TokenPrivateKey)
			},
		},
		{
			name: "token解析失败",
			mock: func(ctrl *gomock.Controller) (svc sms.Service, key []byte) {
				sms := smsmocks.NewMockService(ctrl)
				sms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(jwt.ErrInvalidKey)
				return svc, []byte("1111")
			},
			wantErr: jwt.ErrInvalidKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsSvc, key := tc.mock(ctrl)
			svc := NewSMSService(smsSvc, key)
			err := svc.Send(context.Background(), testToken, []string{"666"}, []string{"18762850585"})
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// @func: TestCreateToken
// @date: 2023-11-10 19:12:20
// @brief: 单元测试-生成token
// @author: Kewin Li
// @param t
func TestCreateToken(t *testing.T) {

	// 设置JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, SMSClaims{
		Tpl: "123456",
	})

	const TokenPrivateKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C5"
	tokenStr, err := token.SignedString([]byte(TokenPrivateKey))
	if err != nil {
		return
	}

	fmt.Printf("token: %s \n", tokenStr)
}
