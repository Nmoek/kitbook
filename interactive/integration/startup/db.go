// Package startup
// @Description: 准备数据库
package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/interactive/repository/dao"
)

func InitDB() *gorm.DB {

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/kitbook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 模块化拆分
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
