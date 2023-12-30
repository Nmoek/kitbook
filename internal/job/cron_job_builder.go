package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"kitbook/pkg/logger"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	l      logger.Logger
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.Logger, opts prometheus.SummaryOpts) *CronJobBuilder {
	vector := prometheus.NewSummaryVec(opts, []string{"job", "success"})

	return &CronJobBuilder{
		l:      l,
		vector: vector,
	}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()

	return cronJobAdapterFunc(func() {
		c.l.DEBUG("job 开始运行", logger.Field{"name", name})

		start := time.Now()

		err := job.Run()
		if err != nil {
			c.l.ERROR("job 执行失败", logger.Error(err), logger.Field{"name", name})
		}

		c.l.DEBUG("job 结束运行", logger.Field{"name", name})

		duration := time.Since(start)
		c.vector.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(duration.Milliseconds()))

	})
}

/*考虑日后接口扩展 */
type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}
