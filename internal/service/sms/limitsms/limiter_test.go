package limitsms

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/service/sms"
	smsmocks "kitbook/internal/service/sms/mocks"
	"kitbook/pkg/limiter"
	limitermocks "kitbook/pkg/limiter/mocks"
	"testing"
)

func TestRateLimitSMSService_Send(t *testing.T) {

	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		// 入参和业务逻辑没有关系可以不进行模拟输入
		ctx     context.Context
		wantErr error
	}{
		{
			name: "验证码发送不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return smsSvc, l
			},
		},
		{
			name: "验证码发送触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)

				return smsSvc, l
			},
			wantErr: ErrIsLimited,
		},
		{
			name: "redis限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("redis限流器错误"))
				return smsSvc, l
			},
			wantErr: errors.New("redis限流器错误"),
		},
		{
			name: "第三方运营商错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				smsSvc := smsmocks.NewMockService(ctrl)
				l := limitermocks.NewMockLimiter(ctrl)

				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				smsSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("腾讯云服务错误"))
				return smsSvc, l
			},
			wantErr: errors.New("腾讯云服务错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsSvc, l := tc.mock(ctrl)
			svc := NewLimitSMSService(smsSvc, l)
			err := svc.Send(tc.ctx, "123456", []string{"666"}, []string{"18762850585"})
			assert.Equal(t, tc.wantErr, err)
		})

	}
}
