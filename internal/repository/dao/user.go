package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

// @func: Insert
// @date: 2023-10-08 02:53:32
// @brief: 数据库插入操作
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param u
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User
// @Description: 表结构定义
type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
}
