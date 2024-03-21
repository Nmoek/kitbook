package dao

import (
	"context"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
)

type FeedPullEventDao interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error

	FindPullEvents(ctx context.Context, ids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error)
	FindPullEventsType(ctx context.Context, typ string, ids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error)
}

type feedPullEventDo struct {
	db *gorm.DB
	l  logger.Logger
}

func NewFeedPullEventDao(db *gorm.DB, l logger.Logger) FeedPullEventDao {
	return &feedPullEventDo{
		db: db,
		l:  l,
	}
}

func (f *feedPullEventDo) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.db.WithContext(ctx).Create(&event).Error
}

func (f *feedPullEventDo) FindPullEvents(ctx context.Context, ids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error) {
	var res []FeedPullEvent
	err := f.db.WithContext(ctx).
		Where("id IN ? AND ctime < ?", ids, timestamp).
		Order("ctime desc").
		Limit(int(limit)).First(&res).Error
	return res, err
}

func (f *feedPullEventDo) FindPullEventsType(ctx context.Context, typ string, ids []int64, timestamp int64, limit int64) ([]FeedPullEvent, error) {
	var res []FeedPullEvent
	err := f.db.WithContext(ctx).Where("id IN ? AND type = ? AND ctime < ?", ids, typ, timestamp).
		Order("ctime desc").
		Limit(int(limit)).First(&res).Error
	return res, err
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
