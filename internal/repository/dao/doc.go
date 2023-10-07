// Package dao
// @Description: 数据库的增删改查
package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
