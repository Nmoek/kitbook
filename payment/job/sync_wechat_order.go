package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"kitbook/payment/service/wechat"
	"kitbook/pkg/logger"
	"sync"
	"time"
)

// SyncWechatOrderJob
// @Description: 定时任务-对账微信订单, 间隔31min
type SyncWechatOrderJob struct {
	svc     *wechat.NativePaymentService
	timeout time.Duration
	client  *rlock.Client
	key     string

	lock      *rlock.Lock
	localLock *sync.Mutex

	l logger.Logger
}

func NewSyncWechatOrderJob(svc *wechat.NativePaymentService, timeout time.Duration, client *rlock.Client, l logger.Logger) *SyncWechatOrderJob {
	return &SyncWechatOrderJob{
		svc:       svc,
		timeout:   timeout,
		client:    client,
		key:       "job_expired_payments",
		localLock: &sync.Mutex{},
		l:         l,
	}
}

func (s *SyncWechatOrderJob) Name() string {
	return "sync_wechat_order_job"
}

func (s *SyncWechatOrderJob) Run() error {

	s.localLock.Lock()
	lock := s.lock
	// 1. 获取分布式锁
	if lock == nil {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		lock, err = s.client.Lock(ctx, s.key, s.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
		cancel()

		if err != nil {
			//TODO: 日志埋点，分布式锁加锁失败
			s.l.WARN("分布式锁加锁失败",
				logger.Error(err),
				logger.Field{"key", s.key})

			// 不返回报错，信任其他web结点会进行处理
			s.localLock.Unlock()
			return nil
		}

		s.lock = lock
		s.localLock.Unlock()

		// 1.1 分布式锁续约
		go func() {
			err2 := s.lock.AutoRefresh(s.timeout/2, s.timeout)
			// 续约失败要将分布式锁置空
			if err2 != nil {
				s.localLock.Lock()
				s.lock = nil
				s.localLock.Unlock()
			}
		}()

	}

	// 2. 找出所有超时过期的订单
	//		2.1 可能有很多, 需要分批找出
	beforeTime := time.Now().Add(-time.Minute * 30)
	offset := 0
	const limit = 100
	// 每次找100条记录, 直到找出的记录不足100条
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		pmts, err := s.svc.FindExpiredPayment(ctx, beforeTime, offset, limit)
		cancel()
		if err != nil {
			return err
		}

		for _, pmt := range pmts {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
			err = s.svc.SyncWechatInfo(ctx, pmt.BizTradeNO)
			cancel()
			if err != nil {
				s.l.ERROR("同步微信订单状态失败",
					logger.Error(err),
					logger.Field{"biz_trade_no", pmt.BizTradeNO})
				continue
			}
		}

		if len(pmts) < limit {
			return nil
		}

		offset += len(pmts)
	}
}
