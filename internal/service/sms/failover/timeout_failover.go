package failover

import (
	"context"
	"kitbook/internal/service/sms"
	"sync/atomic"
)

// FailOverSMSService
// @Description: 装饰器模式-基于超时率选择服务商
type TimeoutFailOverSMSService struct {
	svcs []sms.Service
	// 当前正在使用的服务商
	idx int32
	// 当前服务商已经连续超时了几个请求
	cnt int32
	// 连续超时的阈值
	threshold int32
}

func (t *TimeoutFailOverSMSService) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	length := int32(len(t.svcs))

	if cnt >= t.threshold {
		newIdx := (idx + 1) % length
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 重置连续超时计数
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	err := t.svcs[idx].Send(ctx, templateId, args, phoneNumber)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 遇到错误但不是超时错误
		// 如果强调是超时率. 不增加计数
		// 如果是EOF类的错误，可以选择直接切换服务商
	}

	return err
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}
