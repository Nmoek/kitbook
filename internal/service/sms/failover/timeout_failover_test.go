package failover

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/service/sms"
	smsmocks "kitbook/internal/service/sms/mocks"
	"testing"
)

func TestTimeoutFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		idx       int32
		cnt       int32
		threshold int32

		wantErr error
		wantIdx int32
		wantCnt int32
	}{
		{
			name: "没有触发服务商切换",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc_0 := smsmocks.NewMockService(ctrl)
				svc_0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return []sms.Service{svc_0}
			},
			idx:       0,
			cnt:       12,
			threshold: 15,
			wantErr:   nil,
			wantIdx:   0,
			wantCnt:   0,
		},
		{
			name: "触发服务商切换, 且发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc_0 := smsmocks.NewMockService(ctrl)
				svc_1 := smsmocks.NewMockService(ctrl)
				svc_1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return []sms.Service{svc_0, svc_1}
			},
			idx:       0,
			cnt:       15,
			threshold: 15,
			wantErr:   nil,
			wantIdx:   1,
			wantCnt:   0,
		},
		{
			name: "触发服务商切换, 且发送失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc_0 := smsmocks.NewMockService(ctrl)
				svc_1 := smsmocks.NewMockService(ctrl)
				svc_1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)

				return []sms.Service{svc_0, svc_1}
			},
			idx:       1,
			cnt:       15,
			threshold: 15,
			wantErr:   context.DeadlineExceeded,
			wantIdx:   0,
			wantCnt:   1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsSvc := tc.mock(ctrl)
			svc := NewTimeoutFailOverSMSService(smsSvc, tc.threshold)
			// 特殊情况 要把这几值进行手动预设
			svc.idx = tc.idx
			svc.cnt = tc.cnt
			err := svc.Send(context.Background(), "1223456", []string{"timeout-failover-test"}, []string{"18762850585"})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantIdx, svc.idx)
			assert.Equal(t, tc.wantCnt, svc.cnt)

		})
	}
}
