package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"kitbook/internal/service/sms"
	"time"
)

type Decorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewDecorator(svc sms.Service, opts prometheus.SummaryOpts) *Decorator {
	return &Decorator{

		svc:    svc,
		vector: prometheus.NewSummaryVec(opts, []string{"tpl_id"}),
	}

}

func (d *Decorator) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {
	start := time.Now()

	defer func() {
		duration := time.Since(start).Milliseconds()
		d.vector.WithLabelValues(templateId).Observe(float64(duration))
	}()

	return d.svc.Send(ctx, templateId, args, phoneNumber)
}
