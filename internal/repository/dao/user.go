package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateUser  = errors.New("用户已存在!")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdateById(ctx context.Context, user User) error
	FindByWechat(ctx context.Context, openid string) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func NewGormUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{
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
func (dao *GormUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		// 用户冲突
		if me.Number == duplicateErr {
			// TODO: 当判断手机号时，这里的定义是有问题的
			return ErrDuplicateUser
		}
	}

	return err
}

// @func: FindByEmail
// @date: 2023-10-09 02:02:09
// @brief: 数据库查询操作-按邮箱
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param email
func (dao *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {

	findUser := User{}
	err := dao.db.Where("email = ?", email).First(&findUser).Error
	return findUser, err
}

// @func: FindByID
// @date: 2023-10-12 04:06:06
// @brief: 数据库查询操作-按ID
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param id
// @return error
func (dao *GormUserDao) FindByID(ctx context.Context, id int64) (User, error) {
	findUser := User{}
	err := dao.db.Where("id = ?", id).First(&findUser).Error
	return findUser, err
}

// @func: FindByPhone
// @date: 2023-10-30 00:26:46
// @brief: 数据库查询操作-按手机号
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param phone
// @return interface{}
// @return interface{}
func (dao *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	findUser := User{}
	err := dao.db.Where("phone = ?", phone).First(&findUser).Error
	return findUser, err
}

// @func: UpdateById
// @date: 2023-10-13 01:41:54
// @brief: 数据库更新操作-按ID
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param user
// @return error
func (dao *GormUserDao) UpdateById(ctx context.Context, user User) error {

	return dao.db.WithContext(ctx).Model(&user).Where("id = ?", user.Id).Updates(
		map[string]interface{}{
			"id":       user.Id,
			"nickname": user.Nickname,
			"birthday": user.Birthday,
			"about_me": user.AboutMe,
			"utime":    time.Now().UnixMilli(),
		}).Error
}

// @func: FindByWechat
// @date: 2023-11-12 03:10:44
// @brief: 数据库查询操作-按微信号
// @author: Kewin Li
// @receiver dao
// @param ctx
// @param openid
// @return User
// @return error
func (dao *GormUserDao) FindByWechat(ctx context.Context, openid string) (User, error) {
	findUser := User{}
	err := dao.db.Where("open_id = ?", openid).First(&findUser).Error
	return findUser, err
}

// User
// @Description: 用户表结构定义
type User struct {
	Id    int64          `gorm:"primaryKey, autoIncrement"`
	Email sql.NullString `gorm:"unique"`
	Phone sql.NullString `gorm:"unique"`

	//1. 如果查询要求同时查询openid和unionid, 建立联合唯一索引<openid, unionid>
	//2. 如果查询只要求查询openid，建立一个唯一索引或建立联合唯一索引<openid, unionid>（openid必须在前）
	//3. 如果查询只要求查询unionid，建立一个唯一索引或建立联合唯一索<unionid, openid>
	Openid   sql.NullString `gorm:"unique"`
	Unionid  sql.NullString
	Password string
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`
	Ctime    int64
	Utime    int64
}
