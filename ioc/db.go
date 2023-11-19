// Package ioc
// @Description: 数据库初始化
package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	//dsn := "root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open(viper.GetString("db.dsn")), &gorm.Config{})

	//结构体反序列化, 推荐这种写法
	type Config struct {
		DSN string `yaml:"dsn"`
	}

	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// 初始化表结构(慎该写法)
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
