package opentelemetry

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"kitbook/internal/service/sms"
)

type Decorator struct {
	sms.Service
	tracer trace.Tracer
}

func NewDecorator(svc sms.Service, tracer trace.Tracer) *Decorator {
	return &Decorator{
		Service: svc,
		tracer:  tracer,
	}
}

func (d *Decorator) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	tpCtx, span := d.tracer.Start(ctx, "sms")
	span.SetAttributes(attribute.String("tpl", templateId))
	span.AddEvent("发送短信")
	err := d.Service.Send(tpCtx, templateId, args, phoneNumber)
	if err != nil {
		span.RecordError(err)
	}

	return err
}
