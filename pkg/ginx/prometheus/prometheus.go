package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceId string //实例ID
	Help       string
}

func NewBuilder(namespace string, subsystem string, name string, instanceId string, help string) *Builder {
	return &Builder{
		Namespace:  namespace,
		Subsystem:  subsystem,
		Name:       name,
		InstanceId: instanceId,
		Help:       help,
	}
}

func (b *Builder) BuildResponseTIme() gin.HandlerFunc {
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		// 注意: 以上两个字段都不能有“_”以外的符号
		Name: b.Name + "resp_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)

	prometheus.MustRegister(vector)

	return func(ctx *gin.Context) {
		start := time.Now()

		defer func() {
			// 准备上报prometheus
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()

			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration))
		}()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,

		Name: b.Name + "_active_req",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})

	prometheus.MustRegister(gauge)

	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()

		ctx.Next()
	}

}
