package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

const (
	jobStatusWaiting = iota // 任务等待调度
	jobStatusRunning        // 任务正在运行
	jobStatusPaused         // 任务暂停调度
)

var ErrPreemptJobInvalid = errors.New("抢占任务失败")

type JobDao interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jobId int64) error
	UpdateUtime(ctx context.Context, jobId int64) error
	UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error
}

type GormJobDao struct {
	db *gorm.DB
}

// @func: Preempt
// @date: 2023-12-31 19:22:12
// @brief: MySQL任务调度-任务抢占-乐观锁机制
// @author: Kewin Li
// @receiver g
// @param ctx
// @return Job
// @return error
func (g *GormJobDao) Preempt(ctx context.Context) (Job, error) {
	for {

		var job Job
		now := time.Now().UnixMilli()

		err := g.db.WithContext(ctx).
			Where("next_time < ? AND status = ?", now, jobStatusWaiting).
			First(&job).Error
		if err != nil {
			return job, err
		}

		res := g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", job.Id).Updates(map[string]any{
			"status":  jobStatusRunning,
			"version": job.Version + 1,
			"utime":   now,
		})
		err = res.Error
		if err != nil {
			return Job{}, err
		}

		// 没有抢占到任务
		if res.RowsAffected == 0 {
			continue
		}

		return job, err
	}

}

// @func: Release
// @date: 2023-12-31 19:50:52
// @brief: MySQL任务调度-任务释放
// @author: Kewin Li
// @receiver g
// @param ctx
// @param jobId
// @return error
func (g *GormJobDao) Release(ctx context.Context, jobId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

// @func: UpdateUtime
// @date: 2023-12-31 19:54:10
// @brief: MySQL任务调度-更新utime进行自动续约
// @author: Kewin Li
// @receiver g
// @param ctx
// @param jobId
// @return error
func (g *GormJobDao) UpdateUtime(ctx context.Context, jobId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"utime": now,
	}).Error
}

// @func: UpdateNextTime
// @date: 2023-12-31 21:55:39
// @brief: MySQL任务调度-更新下一次任务调度时间点
// @author: Kewin Li
// @receiver g
// @param ctx
// @param jobId
// @param nextTime
// @return error
func (g *GormJobDao) UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"next_time": nextTime.UnixMilli(),
		"utime":     now,
	}).Error
}

type Job struct {
	Id   int64  `gorm:"primaryKey, autoIncrement"`
	Name string `gorm:"type:varchar(128);unique"`

	Expression string `gorm:"type:varchar(128)"`

	ExecutorName string `gorm:"type:varchar(128)"`
	// 当前任务的抢占状态
	Status int
	// 乐观锁 版本号
	Version int
	// 任务下一次执行的时间点
	NextTime int64 `gorm:"index"`

	Utime int64
	Ctime int64
}
