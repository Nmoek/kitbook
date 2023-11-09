package limitsms

import (
	"context"
	"errors"
	"kitbook/internal/service/sms"
	"kitbook/pkg/limiter"
)

var ErrIsLimited = errors.New("短信验证码发送服务即将被限流")

// LimitSMSService
// @Description: 装饰器模式-短信验证码发送服务
type LimitSMSService struct {
	//	被装饰对象; 核心接口:Send
	svc sms.Service
	// 装饰目的: 执行Send接口前判断限流情况
	limiter limiter.Limiter
	key     string
}

func NewLimitSMSService(svc sms.Service, limiter limiter.Limiter) *LimitSMSService {
	return &LimitSMSService{
		svc:     svc,
		limiter: limiter,
		key:     "sms-limiter",
	}
}

func (r *LimitSMSService) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	isLimited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if isLimited {
		return ErrIsLimited
	}

	return r.svc.Send(ctx, templateId, args, phoneNumber)
}
