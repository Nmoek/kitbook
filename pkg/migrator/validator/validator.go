package validator

import (
	"context"
	"github.com/liyue201/gostl/algorithm/sort"
	"github.com/liyue201/gostl/ds/slice"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
	"kitbook/pkg/migrator"
	"kitbook/pkg/migrator/events"
	"time"
)

type Validator[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	direction string
	batchSize int
	producer  events.Producer
	utime     int64 // 增量校验起始时间点

	sleeInterval time.Duration // >0 睡眠, <=0中断
	fromBase     func(ctx context.Context, offset int) (T, error)
	l            logger.Logger
}

func (v *Validator[T]) Incr() *Validator[T] {
	v.fromBase = v.incrFromBase
	return v
}

func (v *Validator[T]) Full() *Validator[T] {
	v.fromBase = v.fullFromBase
	return v
}

func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var src T
	err := v.base.WithContext(dbCtx).Order("id").Offset(offset).First(&src).Error
	return src, err
}

func (v *Validator[T]) incrFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var src T
	err := v.base.WithContext(dbCtx).
		Where("utime > ?", v.utime).
		Order("utime").
		Offset(offset).
		First(&src).Error
	return src, err
}

// @func: validate
// @date: 2024-01-14 04:02:01
// @brief: 同构数据库数据校验
// @author: Kewin Li
// @receiver v
// @param ctx
// @return error
func (v *Validator[T]) validate(ctx context.Context) error {
	// 同步校验
	//err := v.validateBaseToTarget(ctx)
	//if err != nil {
	//	return err
	//}
	//return v.validateTargetToBase(ctx)

	// 并发校验
	var eg errgroup.Group

	eg.Go(func() error {
		return v.validateBaseToTarget(ctx)
	})

	eg.Go(func() error {
		return v.validateTargetToBase(ctx)
	})

	return eg.Wait()
}

// @func: validateBaseToTarget
// @date: 2024-01-14 03:27:50
// @brief: 正向校验-全量校验/增量校验
// @author: Kewin Li
// @receiver v
// @param ctx
// @return error
func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := -1
	for {
		offset++
		// 1. 从源表查询
		src, err := v.fromBase(ctx, offset)
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			v.l.ERROR("base -> target 查询base失败", logger.Error(err))
			continue
		}

		// 2. 从目标表查询
		var dst T
		err = v.target.WithContext(ctx).Where("id = ?", src.ID()).First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			//发消息到kafka上
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			// 比较两条记录
			isEqual := src.CompareTo(dst)
			if !isEqual {
				//发消息到kafka上
				v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			}
		default:
			v.l.ERROR("base -> target 查询target失败",
				logger.Error(err),
				logger.Int[int64]("id", src.ID()))
			//TODO: 做好protheus数据监控
		}

	}
}

// @func: validateBaseToTargetWithIncrCheck
// @date: 2024-01-14 23:33:48
// @brief: 正向校验-增量校验
// @author: Kewin Li
// @receiver v
// @param ctx
// @return error
func (v *Validator[T]) validateBaseToTargetWithIncrCheck(ctx context.Context) error {
	offset := 0
	for {
		// 1. 从源表查询
		var src T
		err := v.base.WithContext(ctx).
			Where("utime > ?", v.utime).
			Offset(offset).First(&src).Error
		if err == gorm.ErrRecordNotFound {
			if v.sleeInterval <= 0 {
				return nil
			}

			time.Sleep(v.sleeInterval)
			continue
		}
		if err != nil {
			v.l.ERROR("base -> target 查询base失败", logger.Error(err))
			offset++
			continue
		}

		// 2. 从目标表查询
		var dst T
		err = v.target.WithContext(ctx).Where("id = ?", src.ID()).First(&dst).Error
		switch err {
		case gorm.ErrRecordNotFound:
			//发消息到kafka上
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case nil:
			// 比较两条记录
			isEqual := src.CompareTo(dst)
			if !isEqual {
				//发消息到kafka上
				v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			}
		default:
			v.l.ERROR("base -> target 查询target失败",
				logger.Error(err),
				logger.Int[int64]("id", src.ID()))
			//TODO: 做好protheus数据监控
		}
		offset++

	}
}

// @func: validateTargetToBase
// @date: 2024-01-14 03:27:57
// @brief: 反向校验
// @author: Kewin Li
// @receiver v
// @param ctx
// @return error
func (v *Validator[T]) validateTargetToBase(ctx context.Context) error {
	offset := -v.batchSize
	for {
		offset += v.batchSize
		// 1. 查目标表
		var ts []T
		err := v.target.WithContext(ctx).
			Select("id").Order("id").
			Offset(offset).Limit(v.batchSize).Find(&ts).Error

		if err == gorm.ErrRecordNotFound || len(ts) <= 0 {
			return nil
		}
		if err != nil {
			v.l.ERROR("target -> base 查询target失败", logger.Error(err))
			continue
		}

		// 2. 查源表
		var srcTs []T
		ids := make([]int64, len(ts))
		for i, t := range ts {
			ids[i] = t.ID()
		}
		err = v.base.WithContext(ctx).
			Select("id").Where("id IN ?", ids).Find(&srcTs).Error
		if err == gorm.ErrRecordNotFound || len(srcTs) <= 0 {
			v.notifyTargetMissing(ts)
			continue
		}

		if err != nil {
			v.l.ERROR("target -> base, 查询base失败",
				logger.Error(err))
			continue
		}

		// 3. 查出两表差异并修复
		var diffTs []T
		tmpTs := slice.NewSliceWrapper(srcTs)
		for i := 0; i < len(ts) && i < len(srcTs); i++ {
			if !sort.BinarySearch[T](tmpTs.First(), tmpTs.End(), ts[i], func(a, b T) int {
				if a.ID() < b.ID() {
					return -1
				} else if a.ID() > b.ID() {
					return 1
				}

				return 0

			}) {
				diffTs = append(diffTs, ts[i])
			}
		}

		// 4. 所有base没有的记录都需进行kafka消息发送
		v.notifyBaseMissing(diffTs)

		// 5. 批次不够需要退出循环
		if len(ts) < v.batchSize {
			return nil
		}

	}
}

// @func: validateTargetToBaseWithIncrCheck
// @date: 2024-01-14 23:38:51
// @brief: 反向校验-增量校验
// @author: Kewin Li
// @receiver v
// @param ctx
// @return error
func (v *Validator[T]) validateTargetToBaseWithIncrCheck(ctx context.Context) error {
	offset := 0
	for {
		// 1. 查目标表
		var ts []T
		err := v.target.WithContext(ctx).Where("utime = ?", v.utime).
			Select("id").Order("id").
			Offset(offset).Limit(v.batchSize).Find(&ts).Error

		if err == gorm.ErrRecordNotFound || len(ts) <= 0 {
			if v.sleeInterval <= 0 {
				return nil
			}

			time.Sleep(v.sleeInterval)
			continue
		}
		if err != nil {
			v.l.ERROR("target -> base 查询target失败", logger.Error(err))
			offset += v.batchSize
			continue
		}

		// 2. 查源表
		var srcTs []T
		ids := make([]int64, len(ts))
		for i, t := range ts {
			ids[i] = t.ID()
		}
		err = v.base.WithContext(ctx).
			Select("id").Where("id IN ?", ids).Find(&srcTs).Error
		if err == gorm.ErrRecordNotFound || len(srcTs) <= 0 {
			v.notifyTargetMissing(ts)
			offset += v.batchSize

			continue
		}

		if err != nil {
			v.l.ERROR("target -> base, 查询base失败",
				logger.Error(err))
			continue
		}

		// 3. 查出两表差异并修复
		var diffTs []T
		tmpTs := slice.NewSliceWrapper(srcTs)
		for i := 0; i < len(ts) && i < len(srcTs); i++ {
			if !sort.BinarySearch[T](tmpTs.First(), tmpTs.End(), ts[i], func(a, b T) int {
				if a.ID() < b.ID() {
					return -1
				} else if a.ID() > b.ID() {
					return 1
				}

				return 0

			}) {
				diffTs = append(diffTs, ts[i])
			}
		}

		// 4. 所有base没有的记录都需进行kafka消息发送
		v.notifyBaseMissing(diffTs)

		// 5. 批次不够需要退出循环
		if len(ts) < v.batchSize {
			if v.sleeInterval <= 0 {
				return nil
			}

			time.Sleep(v.sleeInterval)

		}

		offset += len(ts)
	}
}

// @func: notify
// @date: 2024-01-14 03:14:27
// @brief: 向kafka发送消息
// @author: Kewin Li
// @receiver v
func (v *Validator[T]) notify(id int64, typ string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		Id:        id,
		Type:      typ,
		Direction: v.direction,
	})

	v.l.ERROR("数据校验不一致消息发送失败",
		logger.Error(err),
		logger.Int[int64]("id", id),
		logger.Field{"type", typ},
		logger.Field{"direction", v.direction})

	//TODO: 接入prothues进行观测
}

// @func: notifyBaseMissing
// @date: 2024-01-14 03:50:28
// @brief: base库数据缺失
// @author: Kewin Li
// @receiver v
func (v *Validator[T]) notifyBaseMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events.InconsistentEventTypeBaseMissing)
	}
}

// @func: notifyTargetMissing
// @date: 2024-01-14 03:51:43
// @brief: target库数据缺失
// @author: Kewin Li
// @receiver v
// @param ts
func (v *Validator[T]) notifyTargetMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events.InconsistentEventTypeTargetMissing)
	}
}
