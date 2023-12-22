package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"kitbook/internal/domain"
	"kitbook/internal/service/oauth2/wechat"
)

type Decorator struct {
	wechat.Service
	tracer trace.Tracer
}

func NewDecorator(svc wechat.Service, tracer trace.Tracer) *Decorator {
	return &Decorator{
		Service: svc,
		tracer:  tracer,
	}
}

func (d  *Decorator) VerifyCode(ctx context.Context, code string) (domain.WechtInfo, error)
	ctx, span := d.tracer.Start(ctx, "wechat")
	span.AddEvent("微信验证登录")
	err := d.Service.VerifyCode(ctx, code)
	if err != nil {
		span.RecordError(err)
	}

	return err
}
