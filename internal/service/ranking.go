package service

import (
	"context"
	"github.com/liyue201/gostl/ds/priorityqueue"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	intrSvc InteractiveService
	artSvc  ArticleService

	repo repository.RankingRepository

	batchSize int //一批查询出多少条数据
	// 分数生成函数
	scoreFunc func(likeCnt int64, utime time.Time) float64
	n         int //总共需要多少条数据
}

func NewBatchRankingService(intrSvc InteractiveService, artSvc ArticleService) RankingService {
	return &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		batchSize: 100, // 每一批查100条记录
		n:         100, // 维护score最高的前100条记录
		scoreFunc: func(likeCnt int64, utime time.Time) float64 {
			duration := time.Since(utime).Seconds()
			return float64(likeCnt-1) / math.Pow(duration+1, 1.5)
		},
	}
}

// @func: TopN
// @date: 2023-12-27 00:27:50
// @brief: 热点算法
// @author: Kewin Li
// @receiver b
// @param ctx
// @return error
func (b *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := b.topN(ctx)
	if err != nil {
		return err
	}

	//最终放入缓存中
	//拆分两个函数, 目的：先将热点算法本身测试正确，然后再连带缓存放入一起进行测试
	return b.repo.ReplaceTopN(ctx, arts)
}

// @func: topN
// @date: 2023-12-27 00:27:42
// @brief: 热点算法真正实现
// @author: Kewin Li
// @receiver b
// @param ctx
// @return error
func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	type Score struct {
		score float64
		art   domain.Article
	}

	offset := 0
	start := time.Now()
	// 截止一周前数据
	ddl := start.Add(-7 * 24 * time.Hour)

	queue := priorityqueue.New[Score](func(a, b Score) int {
		if a.score < b.score {
			return -1

		} else if a.score > b.score {
			return 1
		}

		if a.art.Utime.UnixMilli() < b.art.Utime.UnixMilli() {
			return -1
		} else if a.art.Utime.UnixMilli() > b.art.Utime.UnixMilli() {
			return 1
		}

		return 0
	})

	for {
		// 查文章
		arts, err := b.artSvc.ListPub(ctx, start, offset, b.batchSize)
		if err != nil {
			return nil, err
		}

		// 获取的帖子数据为空需要退出
		if len(arts) <= 0 {
			//TODO: 日志埋点, 哪一次获取的数据为空
			break
		}

		bizIds := make([]int64, 0, len(arts))
		for _, art := range arts {
			bizIds = append(bizIds, art.Id)
		}

		// 取出点赞数
		intrMap, err := b.intrSvc.GetByIds(ctx, "article", bizIds)
		if err != nil {
			return nil, err
		}

		// 计算score
		for _, art := range arts {

			score := b.scoreFunc(intrMap[art.Id].LikeCnt, art.Utime)

			// 队列达到N后
			if queue.Size() >= b.n {
				top := queue.Top()
				// 当前计算的score是否比堆顶score更小
				// 是，不放入队列
				// 否，拿出堆顶，进行更新
				if score < top.score {
					continue
				}

				queue.Pop()
			}

			// 放入优先队列
			queue.Push(Score{
				score: score,
				art:   art,
			})
		}

		offset += len(arts)
		// 当前这一批数据已经不满足取出的阈值
		// 或  当前这一批数据的最后一条记录的utime已经超出了一周前
		// 认为没有下一批符合的数据, 就直接退出
		if len(arts) < b.batchSize || arts[len(arts)-1].Utime.Before(ddl) {
			break
		}
	}

	//将数据从小顶堆中取出
	res := make([]domain.Article, queue.Size())
	for i := queue.Size() - 1; i >= 0 && !queue.Empty(); i-- {
		res[i] = queue.Top().art
		queue.Pop()
	}

	return res, nil

}

// @func: GetTopN
// @date: 2023-12-30 21:45:39
// @brief: 热榜服务-查询热榜数据
// @author: Kewin Li
// @receiver b
// @param ctx
// @return error
func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return b.repo.GetTopN(ctx)
}
