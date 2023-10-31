// Package ioc
// @Description: 数据库初始化
package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/config"
	"kitbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
