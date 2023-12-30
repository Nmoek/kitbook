package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"kitbook/internal/service"
	"kitbook/pkg/logger"
	"sync"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
	client  *rlock.Client
	key     string

	lock      *rlock.Lock
	localLock *sync.Mutex

	l logger.Logger
}

func NewRankingJob(svc service.RankingService,
	timeout time.Duration,
	client *rlock.Client,
	l logger.Logger) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		key:       "job_ranking",
		client:    client,
		localLock: &sync.Mutex{},
		l:         l,
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	var err error
	r.localLock.Lock()
	lock := r.lock

	if lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()

		lock, err = r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
		if err != nil {
			//TODO: 日志埋点，分布式锁加锁失败
			r.l.WARN("分布式锁加锁失败", logger.Error(err))
			// 不返回报错，信任其他web结点会进行处理
			r.localLock.Unlock()
			return nil
		}

		r.lock = lock
		r.localLock.Unlock()

		// 重要：续约机制
		go func() {
			// 间隔多久续约一次?  续约的时长是多长?
			err2 := r.lock.AutoRefresh(r.timeout/2, r.timeout)
			if err2 != nil {
				//TODO: 日志埋点，自动续约出错

				r.localLock.Lock()
				r.lock = nil
				r.localLock.Unlock()

			}
		}()
	}

	// 执行业务
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

// @func: Run
// @date: 2023-12-31 00:06:31
// @brief: 分布式锁-朴素加解锁方案
// @author: Kewin Li
// @receiver r
// @return error
func (r *RankingJob) RunV1() error {

	//1. 分布式锁 加锁
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	// 整个加锁过程，持有锁时间=热榜计算超时时间
	lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
		Interval: time.Millisecond * 100, //每次重试加锁间隔 100ms
		Max:      3,                      //重试加锁次数 3次
	}, time.Second) // 重试加锁总共超时 1s

	if err != nil {
		return err
	}

	// 3. 任务执行完毕, 分布式锁进行解锁
	defer func() {
		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
		defer cancel2()

		//注意: 解锁失败也没有必要进行重试，最多r.timeout的时间也会自动解锁
		err2 := lock.Unlock(ctx2)
		if err2 != nil {
			//TODO: 日志埋点，分布式锁解锁失败
			r.l.ERROR("分布式锁解锁失败",
				logger.Error(err),
				logger.Field{"name", r.Name()})
		}
	}()

	// 2. 执行定时任务
	ctx, cancel = context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}
