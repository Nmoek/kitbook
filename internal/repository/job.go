package repository

import (
	"context"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
	"time"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jobId int64) error
	UpdateUtime(ctx context.Context, jobId int64) error
	UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error
}

type PreemptJobRepository struct {
	dao dao.JobDao
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id:           job.Id,
		Expression:   job.Expression,
		Name:         job.Name,
		ExecutorName: job.ExecutorName,
	}, err

}

func (p *PreemptJobRepository) Release(ctx context.Context, jobId int64) error {
	return p.dao.Release(ctx, jobId)
}

// @func: UpdateUtime
// @date: 2023-12-31 19:49:56
// @brief: MySQL任务调度-更新续约时间
// @author: Kewin Li
// @receiver p
// @param jobId
// @return error
func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, jobId int64) error {
	return p.dao.UpdateUtime(ctx, jobId)
}

// @func: UpdateNextTime
// @date: 2023-12-31 21:54:37
// @brief: MySQL任务调度-更新下一次任务调度时间点
// @author: Kewin Li
// @receiver p
// @param ctx
// @param jobId
// @param nextTime
// @return error
func (p *PreemptJobRepository) UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error {
	return p.dao.UpdateNextTime(ctx, jobId, nextTime)
}
