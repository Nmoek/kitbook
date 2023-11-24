package dao

import (
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},    //用户表
		&Article{}, //帖子表
	)
}
