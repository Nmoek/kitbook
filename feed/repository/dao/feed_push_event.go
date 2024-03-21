package dao

import (
	"context"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
)

type FeedPushEventDao interface {
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	FindPushEvents(ctx context.Context, id int64, timestamp int64, limit int64) ([]FeedPushEvent, error)
	FindPushEventsType(ctx context.Context, typ string, id int64, timestamp int64, limit int64) ([]FeedPushEvent, error)
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

func (f *feedPushEventDao) FindPushEvents(ctx context.Context, id int64, timestamp int64, limit int64) ([]FeedPushEvent, error) {
	var res []FeedPushEvent
	err := f.db.WithContext(ctx).Where("id = ? AND ctime < ?", id, timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (f *feedPushEventDao) FindPushEventsType(ctx context.Context, typ string, id int64, timestamp int64, limit int64) ([]FeedPushEvent, error) {
	var res []FeedPushEvent
	err := f.db.WithContext(ctx).
		Where("id = ? AND type = ? AND ctime < ?", id, typ, timestamp).
		Order("ctime decs").
		Limit(int(limit)).Find(&res).Error
	return res, err
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
