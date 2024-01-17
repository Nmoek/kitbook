package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"kitbook/internal/domain"
	"kitbook/internal/service/oauth2/wechat"
	"time"
)

type Decorator struct {
	wechat.Service
	sum prometheus.Summary
}

func NewDecorator(svc wechat.Service, sum prometheus.Summary) *Decorator {
	return &Decorator{
		Service: svc,
		sum:     sum,
	}
}

func (d *Decorator) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Microseconds()
		d.sum.Observe(float64(duration))
	}()

	return d.Service.VerifyCode(ctx, code)
}
