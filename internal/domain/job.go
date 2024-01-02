package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id           int64
	Name         string
	Expression   string
	ExecutorName string
	CancelFunc   func()
}

// @func: NextTime
// @date: 2023-12-31 21:52:02
// @brief: 利用cron表达式去计算下一次任务的调度时间点
// @author: Kewin Li
// @receiver j
// @return time.Time
func (j Job) NextTime() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}
