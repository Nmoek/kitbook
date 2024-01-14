package dao

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
)

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

var ErrUnkonwPattern = errors.New("未知pattern")

type DoubleWriteDAO struct {
	src InteractiveDao
	dst InteractiveDao

	pattern *atomicx.Value[string]

	l logger.Logger
}

func NewDoubleWriteDAO(src *gorm.DB, dst *gorm.DB, l logger.Logger) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     NewGORMInteractiveDao(src),
		dst:     NewGORMInteractiveDao(dst),
		pattern: atomicx.NewValueOf(PatternSrcOnly),
		l:       l,
	}
}

// @func: UpdatePattern
// @date: 2024-01-14 18:14:33
// @brief: 模式切换
// @author: Kewin Li
// @receiver d
// @param pattern
func (d *DoubleWriteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}

// TODO: 双写 写入演示
func (d *DoubleWriteDAO) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		return d.src.IncreaseReadCnt(ctx, biz, bizId)
	case PatternSrcFirst:
		err := d.src.IncreaseReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}
		err = d.dst.IncreaseReadCnt(ctx, biz, bizId)
		if err != nil {
			// 不进行return, 正常来说双写阶段, src成功就算业务上成功了
			d.l.ERROR("双写写入dst失败",
				logger.Error(err),
				logger.Field{"biz", biz},
				logger.Field{"biz_id", bizId})
		}
		return nil

	case PatternDstFirst:
		err := d.dst.IncreaseReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}
		err = d.src.IncreaseReadCnt(ctx, biz, bizId)
		if err != nil {
			// 不进行return, 正常来说双写阶段, dst成功就算业务上成功了
			d.l.ERROR("双写写入src失败",
				logger.Error(err),
				logger.Field{"biz", biz},
				logger.Field{"biz_id", bizId})
		}

		return nil
	case PatternDstOnly:
		return d.dst.IncreaseReadCnt(ctx, biz, bizId)

	default:
		return ErrUnkonwPattern
	}

}

func (d *DoubleWriteDAO) AddLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) DelLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) (UserLikeInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) AddCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) DelCollectionItem(ctx context.Context, biz string, bizId int64, collectId int64, userId int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetCollectionItem(ctx context.Context, biz string, bizId int64, userId int64) (UserCollectInfo, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: 双写 读取演示
func (d *DoubleWriteDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.Get(ctx, biz, bizId)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.Get(ctx, biz, bizId)
	default:
		return Interactive{}, ErrUnkonwPattern
	}
}

func (d *DoubleWriteDAO) GetV1(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		intr, err := d.src.Get(ctx, biz, bizId)
		if err == nil {
			// 并发进行一次增量校验
			// 缺点是校验过程耦合在了查询过程中
			go func() {
				intrDst, err2 := d.dst.Get(ctx, biz, bizId)
				if err2 != nil {
					if intr != intrDst {
						d.l.ERROR("查询过程中数据校验不一致",
							logger.Field{"biz", biz},
							logger.Field{"biz_id", bizId})
					}
					d.l.ERROR("增量校验是dst查询失败",
						logger.Error(err2),
						logger.Field{"biz", biz},
						logger.Field{"biz_id", bizId})
				}
			}()
		}
		return intr, err
	case PatternDstOnly, PatternDstFirst:
		return d.dst.Get(ctx, biz, bizId)
	default:
		return Interactive{}, ErrUnkonwPattern
	}
}

func (d *DoubleWriteDAO) BatchIncreaseReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetByIds(ctx context.Context, biz string, bizIds []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}
