package connpool

import (
	"context"
	"database/sql"
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

// DoubleWritePool
// @Description: 双写连接池
type DoubleWritePool struct {
	src gorm.ConnPool
	dst gorm.ConnPool

	pattern *atomicx.Value[string]

	l logger.Logger
}

func NewDoubleWritePool(src *gorm.DB, dst *gorm.DB, l logger.Logger) *DoubleWritePool {
	return &DoubleWritePool{
		src:     src.ConnPool,
		dst:     dst.ConnPool,
		pattern: atomicx.NewValueOf(PatternSrcOnly),
		l:       l,
	}
}

func (d *DoubleWritePool) UpdatePattern(pattern string) error {
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst, PatternDstFirst, PatternDstOnly:
		d.pattern.Store(pattern)
	default:
		return ErrUnkonwPattern
	}
	return nil
}

func (d DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pateern := d.pattern.Load()
	switch pateern {
	case PatternSrcOnly:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{src: src, l: d.l, pattern: pateern}, err
	case PatternSrcFirst:
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			d.l.ERROR("双写事务dst表操作失败", logger.Error(err))
		}
		return &DoubleWriteTx{src: src, dst: dst, l: d.l, pattern: pateern}, nil
	case PatternDstFirst:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		src, err := d.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			d.l.ERROR("双写事务src表操作失败", logger.Error(err))
		}
		return &DoubleWriteTx{src: src, dst: dst, l: d.l, pattern: pateern}, nil
	case PatternDstOnly:
		dst, err := d.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWriteTx{dst: dst}, err
	default:
		return nil, ErrUnkonwPattern
	}
}

func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO 这个方法无法进行装饰
	panic("double write dont support")
}

func (d *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写src失败", logger.Error(err), logger.Field{"sql", query})
			return res, err
		}

		res, err = d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写dst失败", logger.Error(err), logger.Field{"sql", query})
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写dst失败", logger.Error(err), logger.Field{"sql", query})
			return res, err
		}

		res, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写src失败", logger.Error(err), logger.Field{"sql", query})
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)

	default:
		return nil, ErrUnkonwPattern
	}
}

func (d *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, ErrUnkonwPattern
	}

}

func (d *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		//TODO: 注意 这里是天坑
		// 1. 使用unsafe强行操作内存赋值
		// 2. 直接panic, 但需要业务代码中去处理这个panic
		panic(ErrUnkonwPattern)
	}
}

// DoubleWriteTx
// @Description: 双写事务
type DoubleWriteTx struct {
	src *sql.Tx
	dst *sql.Tx
	//注意：某个事务卡住了很久都没有结束造成的数据不一致问题是无法解决的, 最终依赖于数据校验与修复流程
	pattern string
	l       logger.Logger
}

func (d *DoubleWriteTx) Commit() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Commit()
	case PatternSrcFirst:
		err := d.src.Commit()
		// src提交失败该怎么办?
		if err != nil {
			return err
		}

		if d.dst != nil {
			err = d.dst.Commit()
			if err != nil {
				//不能返回, 打印日志, 通过数据校验去修复
				d.l.ERROR("dst表提交事务失败", logger.Error(err))
			}
		}
		return nil
	case PatternDstFirst:
		err := d.dst.Commit()
		// dst提交失败该怎么办?
		if err != nil {
			return err
		}

		if d.src != nil {
			err = d.src.Commit()
			if err != nil {
				//不能返回, 打印日志, 通过数据校验去修复
				d.l.ERROR("src表提交事务失败", logger.Error(err))
			}
		}
		return nil
	case PatternDstOnly:
		return d.dst.Commit()
	default:
		return ErrUnkonwPattern
	}
}

func (d *DoubleWriteTx) Rollback() error {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.Rollback()
	case PatternSrcFirst:
		err := d.src.Rollback()
		if err != nil {
			return err
		}

		err = d.dst.Rollback()
		if err != nil {
			d.l.ERROR("src表回滚失败", logger.Error(err))
		}
		return nil
	case PatternDstFirst:
		err := d.src.Rollback()
		if err != nil {
			return err
		}

		err = d.dst.Rollback()
		if err != nil {
			d.l.ERROR("src表回滚失败", logger.Error(err))
		}
		return nil
	case PatternDstOnly:
		return d.dst.Rollback()
	default:
		return ErrUnkonwPattern
	}

}

func (d *DoubleWriteTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	//TODO 这个方法无法进行装饰
	panic("double write dont support")
}

func (d *DoubleWriteTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch d.pattern {
	case PatternSrcOnly:
		return d.src.ExecContext(ctx, query, args...)
	case PatternSrcFirst:
		res, err := d.src.ExecContext(ctx, query, args...)
		if err != nil && d.dst != nil {
			d.l.ERROR("双写src失败", logger.Error(err), logger.Field{"sql", query})
			return res, err
		}

		res, err = d.dst.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写dst失败", logger.Error(err), logger.Field{"sql", query})
		}
		return res, err
	case PatternDstFirst:
		res, err := d.dst.ExecContext(ctx, query, args...)
		//d.src != nil代表事务开启成功
		if err != nil && d.src != nil {
			d.l.ERROR("双写dst失败", logger.Error(err), logger.Field{"sql", query})
			return res, err
		}

		res, err = d.src.ExecContext(ctx, query, args...)
		if err != nil {
			d.l.ERROR("双写src失败", logger.Error(err), logger.Field{"sql", query})
		}
		return res, err
	case PatternDstOnly:
		return d.dst.ExecContext(ctx, query, args...)

	default:
		return nil, ErrUnkonwPattern
	}
}

func (d *DoubleWriteTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryContext(ctx, query, args...)
	default:
		return nil, ErrUnkonwPattern
	}

}

func (d *DoubleWriteTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch d.pattern {
	case PatternSrcOnly, PatternSrcFirst:
		return d.src.QueryRowContext(ctx, query, args...)
	case PatternDstFirst, PatternDstOnly:
		return d.dst.QueryRowContext(ctx, query, args...)
	default:
		//TODO: 注意 这里是天坑
		// 1. 使用unsafe强行操作内存赋值
		// 2. 直接panic, 但需要业务代码中去处理这个panic
		panic(ErrUnkonwPattern)
	}
}
