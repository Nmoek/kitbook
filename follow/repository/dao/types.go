package dao

import "context"

type FollowDao interface {
	CreateFollowRelation(ctx context.Context, f FollowRelation) error
	UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error
	FollowRelationList(ctx context.Context, follower int64, offset int64, limit int64) ([]FollowRelation, error)
	FollowRelationDetail(ctx context.Context, followee int64, follower int64) (FollowRelation, error)
	CntFollower(ctx context.Context, uid int64) (int64, error)
	CntFollowee(ctx context.Context, uid int64) (int64, error)
}

// FollowRelation
// @Description: 用户关注数据表
type FollowRelation struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Followee int64 `gorm:"uniqueIndex:follower_followee"`
	Follower int64 `gorm:"uniqueIndex:follower_followee"`

	Status uint8

	//TODO: 关注中还有细分类型、优先级、用户分组等等
	//Type string
	//Priority string
	//GroupId int64

	Utime int64
	Ctime int64
}

const (
	FollowRelationStatusUnknown uint8 = iota
	FollowRelationStatusActive
	FollowRelationStatusInactive
)
