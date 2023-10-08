package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrDuplicateEmail = errors.New("邮箱已经被注册!")

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
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		// 用户冲突
		if me.Number == duplicateErr {
			// TODO: 当判断手机号时，这里的定义是有问题的
			return ErrDuplicateEmail
		}
	}

	return err
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
