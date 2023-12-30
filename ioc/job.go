package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"kitbook/internal/job"
	"kitbook/internal/service"
	"kitbook/pkg/logger"
	"time"
)

func InitRankingJob(svc service.RankingService, client *rlock.Client, l logger.Logger) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30, client, l)
}

func InitJobs(l logger.Logger, ranking_job *job.RankingJob) *cron.Cron {

	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "kewin",
		Subsystem: "kitbook",
		Name:      "corn_job",
		Help:      "定时任务执行",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m", builder.Build(ranking_job))
	if err != nil {
		panic(err)
	}

	return expr
}
