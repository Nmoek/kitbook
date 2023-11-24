// Package startup
// @Description: 准备数据库
package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/internal/repository/dao"
)

func InitDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
