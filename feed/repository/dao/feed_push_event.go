package dao

import (
	"context"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
)

type FeedPushEventDao interface {
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
}

type feedPushEventDao struct {
	db *gorm.DB
	l  logger.Logger
}

func NewFeedPushEventDao(db *gorm.DB, l logger.Logger) FeedPushEventDao {
	return &feedPushEventDao{
		db: db,
		l:  l,
	}
}

func (f *feedPushEventDao) CreatePushEvents(ctx context.Context, events []FeedPushEvent) error {
	return f.db.WithContext(ctx).Create(&events).Error
}

// FeedPushEvent
// @Description: 推事件(对应收件箱)
type FeedPushEvent struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	// 收件人
	Uid  int64  `gorm:"index"`
	Type string `gorm:"type:varchar(256)"`
	// 该字段存放拓展字段
	Content string `gorm:"type:text"`
	Ctime   int64
}
