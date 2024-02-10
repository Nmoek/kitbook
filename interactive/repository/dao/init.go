package dao

import (
	"gorm.io/gorm"
	"kitbook/reward/repository/dao"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&dao.Reward{},
	)
}
