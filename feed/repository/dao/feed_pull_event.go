package dao

import (
	"context"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
)

type FeedPullEventDao interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
}

type feedPullEventDap struct {
	db *gorm.DB
	l  logger.Logger
}

func NewFeedPullEventDao(db *gorm.DB, l logger.Logger) FeedPullEventDao {
	return &feedPullEventDap{
		db: db,
		l:  l,
	}
}

func (f *feedPullEventDap) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}

// FeedPullEvent
// @Description: 拉事件(对应发件箱)
type FeedPullEvent struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	// 收件人
	Uid  int64  `gorm:"index"`
	Type string `gorm:"type:varchar(256)"`
	// 该字段存放拓展字段
	Content string `gorm:"type:text"`
	Ctime   int64
}
