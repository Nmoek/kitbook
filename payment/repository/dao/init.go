package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Payment{}, //用户收藏信息表
	)
}
