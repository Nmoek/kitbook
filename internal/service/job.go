package service

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/pkg/logger"
	"time"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
}

type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.Logger
}

func NewCronJobService(repo repository.JobRepository, l logger.Logger) *CronJobService {
	return &CronJobService{
		repo:            repo,
		refreshInterval: time.Minute,
		l:               l}
}

// @func: Preempt
// @date: 2023-12-31 19:10:19
// @brief: MySQL任务调度-抢占任务
// @author: Kewin Li
// @receiver c
// @param ctx
// @return domain.Job
// @return error
func (c *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}

	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		// 等待定时器到期
		for range ticker.C {
			c.refresh(ctx, job.Id)
		}
	}()

	job.CancelFunc = func() {
		// 关闭自动续约，否则goroutine泄漏
		ticker.Stop()
		ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err2 := c.repo.Release(ctx2, job.Id)
		if err2 != nil {
			c.l.ERROR("释放定时任务失败",
				logger.Error(err2),
				logger.Int[int64]("job_id", job.Id))
		}
	}

	return job, nil
}

// //  @func: refresh
//
//	@date: 2023-12-31 19:46:13
//	@brief: MySQL任务调度-抢占时间续约
//	@author: Kewin Li
//	@receiver c
//	@param jobId
func (c *CronJobService) refresh(ctx context.Context, jobId int64) {
	err := c.repo.UpdateUtime(ctx, jobId)
	if err != nil {
		c.l.ERROR("任务续约失败",
			logger.Error(err),
			logger.Int[int64]("jog_id", jobId))
	}
}

// @func: ResetNextTime
// @date: 2023-12-31 21:45:15
// @brief: MySQL任务调度-任务执行完毕后重新计算下一次调度时间点
// @author: Kewin Li
// @receiver c
// @param ctx
// @param job
// @return error
func (c *CronJobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	nextTime := job.NextTime()
	return c.repo.UpdateNextTime(ctx, job.Id, nextTime)
}
