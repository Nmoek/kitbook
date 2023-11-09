package failover

import (
	"context"
	"errors"
	"kitbook/internal/service/sms"
	"sync/atomic"
)

var ErrAllServiceSendFail = errors.New("所有服务商都发送失败")

// FailOverSMSServiceV1
// @Description: 装饰器模式-轮询服务商模式
type FailOverSMSServiceV1 struct {
	svcs []sms.Service
}

func (f *FailOverSMSServiceV1) SendV1(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, templateId, args, phoneNumber)
		if err == nil {
			return nil
		}

		// TODO: 日志埋点
		//fmt.Printf("svc info, err! %s \n", err)

	}

	return ErrAllServiceSendFail
}

// FailOverSMSServiceV2
// @Description: 装饰器模式-递增当前选择的服务商的下标
type FailOverSMSService struct {
	svcs []sms.Service
	// 上一次使用的服务商
	idx uint64
}

func (f *FailOverSMSService) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	// TODO: 这种写法有两个问题:1.CPU高速缓存带来的值不统一;2.多线程读写问题
	//idx := f.idx
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; idx++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, templateId, args, phoneNumber)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			//TODO:日志埋点,服务商调用失败、超时、取消等
			return err
		}

	}

	return ErrAllServiceSendFail
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
		idx:  0,
	}
}
