package job

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	"kitbook/pkg/logger"
	"time"
)

// Executor
// @Description: 执行器
type Executor interface {
	Name() string
	// 这里的ctx是全局控制, 特别注意：Executor的实现需要关注ctx的超时、取消
	Exec(ctx context.Context, job domain.Job) error
}

// LocalFuncExecutor
// @Description: 本地任务调用
type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, job domain.Job) error
}

func NewLocalFuncExecutor(funcs map[string]func(ctx context.Context, job domain.Job) error) *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: funcs,
	}
}

// @func: RegisterFunc
// @date: 2023-12-31 22:11:38
// @brief: 执行器-任务注册
// @author: Kewin Li
// @receiver l
func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, job domain.Job) error) {
	l.funcs[name] = fn
}

// @func: Name
// @date: 2023-12-31 22:09:35
// @brief: 执行器-标识本地任务、远程任务
// @author: Kewin Li
// @receiver l
// @return string
func (l *LocalFuncExecutor) Name() string {
	return "local"
}

// @func: Exec
// @date: 2023-12-31 22:09:23
// @brief: 执行器-任务真正执行
// @author: Kewin Li
// @receiver l
// @param ctx
// @param job
// @return error
func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.Job) error {
	fn, ok := l.funcs[job.Name]
	if !ok {
		return fmt.Errorf("未注册的本地方法: [%d] %s", job.Id, job.Name)
	}

	return fn(ctx, job)
}

// Scheduler
// @Description: 调度器
type Scheduler struct {
	svc service.JobService

	dbTimeout time.Duration
	executors map[string]Executor

	// 令牌算法进行限流
	limiter *semaphore.Weighted

	l logger.Logger
}

func NewScheduler(svc service.JobService, l logger.Logger) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		executors: map[string]Executor{},
		limiter:   semaphore.NewWeighted(100), //同一个web实例最多同时运行100个任务
		l:         l}
}

// @func: RegisterExecutor
// @date: 2023-12-31 22:11:06
// @brief: 调度器-执行器注册
// @author: Kewin Li
// @receiver s
// @param exec
func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

// @func: Schedule
// @date: 2023-12-31 21:30:46
// @brief: 调度器-定时任务调度
// @author: Kewin Li
// @receiver s
func (s *Scheduler) Schedule(ctx context.Context) {
	for {

		/*任务抢占保护 start*/
		if ctx.Err() != nil {
			s.l.INFO("context 出错/超时", logger.Error(ctx.Err()))
			return
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			s.l.WARN("令牌获取失败", logger.Error(err))
			return
		}
		/*任务抢占保护 end*/

		dbCtx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)

		// 1. 抢占任务
		job, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			s.l.WARN("任务调度失败",
				logger.Error(err),
				logger.Int[int64]("job_id", job.Id))
			continue
		}

		// 2. 执行任务
		exec, ok := s.executors[job.ExecutorName]
		// 2.1 没有发现执行器
		if !ok {
			s.l.ERROR("执行器未发现",
				logger.Field{"ok", ok},
				logger.Int[int64]("job_id", job.Id),
				logger.Field{"executor_name", job.ExecutorName})
		}

		// 2.2 任务开始执行
		go func() {
			defer func() {
				s.limiter.Release(1) //释放令牌

				job.CancelFunc() // 资源释放
			}()

			err2 := exec.Exec(ctx, job)
			if err2 != nil {
				s.l.ERROR("任务执行发生错误",
					logger.Error(err2),
					logger.Int[int64]("job_id", job.Id),
					logger.Field{"executor_name", job.ExecutorName})
			}

			// 2.3 执行完毕，更新下一次任务调度的时间点
			err2 = s.svc.ResetNextTime(ctx, job)

		}()
	}
}
