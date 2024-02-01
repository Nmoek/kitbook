package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"kitbook/internal/job"
	job2 "kitbook/payment/job"
	"kitbook/payment/service/wechat"
	"kitbook/pkg/logger"
	"time"
)

func InitSyncWechatOrderJob(svc *wechat.NativePaymentService, client *rlock.Client, l logger.Logger) *job2.SyncWechatOrderJob {
	return job2.NewSyncWechatOrderJob(svc, time.Second*30, client, l)
}

func InitJobs(l logger.Logger, wechat_job *job2.SyncWechatOrderJob) *cron.Cron {

	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "kewin",
		Subsystem: "payment",
		Name:      "corn_job",
		Help:      "批量处理超时订单定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m", builder.Build(wechat_job))
	if err != nil {
		panic(err)
	}

	return expr
}
