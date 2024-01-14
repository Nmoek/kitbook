package fixer

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"kitbook/pkg/migrator"
	"kitbook/pkg/migrator/events"
)

var ErrUnknowEventType = errors.New("未知事件类型")

type OverrideFixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB

	columns []string
}

func NewOverrideFixer[T migrator.Entity](base *gorm.DB, target *gorm.DB) (*OverrideFixer[T], error) {
	rows, err := base.Model(new(T)).Order("id").Rows()
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()

	return &OverrideFixer[T]{
		base:    base,
		target:  target,
		columns: columns,
	}, nil
}

func (f *OverrideFixer[T]) Fix(event events.InconsistentEvent) error {
	switch event.Type {
	case events.InconsistentEventTypeNEQ,
		events.InconsistentEventTypeTargetMissing:

		var t T
		err := f.base.Where("id = ?", event.Id).First(&t).Error
		switch err {
		// base找不到记录是因为消息的处理可能非常的滞后
		case gorm.ErrRecordNotFound:
			return f.target.Model(&t).Delete("id = ?", event.Id).Error
		case nil:
			// 注意：upsert语义解决并发问题
			//return f.target.Model(&t).Updates(&t).Error
			return f.target.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Error
		default:
			return err
		}

	case events.InconsistentEventTypeBaseMissing:
		return f.target.Model(new(T)).Delete("id = ?", event.Id).Error

	default:
		return ErrUnknowEventType
	}

}

func (f *OverrideFixer[T]) FixV1(ctx context.Context, id int64) error {
	// 粗暴做法
	var t T
	err := f.base.WithContext(ctx).Where("id = ?", id).First(&t).Error
	switch err {
	// base找不到记录是因为消息的处理可能非常的滞后
	case gorm.ErrRecordNotFound:
		return f.target.WithContext(ctx).Model(&t).Delete("id = ?", id).Error
	case nil:
		// 注意：upsert语义解决并发问题
		//return f.target.Model(&t).Updates(&t).Error
		return f.target.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Error
	default:
		return err
	}

}
