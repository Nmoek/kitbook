package failover

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/service/sms"
	smsmocks "kitbook/internal/service/sms/mocks"
	"testing"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		wantErr error
	}{
		{
			name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc_0 := smsmocks.NewMockService(ctrl)

				svc_0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

				return []sms.Service{svc_0}
			},
			wantErr: nil,
		},
		{
			name: "部分发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svcs := make([]sms.Service, 10)
				for i := 0; i < 10; i++ {
					s := smsmocks.NewMockService(ctrl)

					if i == 9 {
						s.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
					} else {
						s.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(fmt.Sprintf("%d 服务商错误", i)))

					}
					svcs[i] = s
				}

				return svcs
			},
			wantErr: nil,
		},
		{
			name: "全部发送失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svcs := make([]sms.Service, 10)
				for i := 0; i < 10; i++ {
					s := smsmocks.NewMockService(ctrl)
					s.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(fmt.Sprintf("%d 服务商错误", i)))
					svcs[i] = s
				}

				return svcs
			},
			wantErr: ErrAllServiceSendFail,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewFailOverSMSService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "123456", []string{"failover-test"}, []string{"18762850585"})
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
