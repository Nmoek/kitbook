// Package ioc
// @Description: 数据库初始化
package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"kitbook/internal/repository/dao"
	"kitbook/pkg/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
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
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.DEBUG), glogger.Config{
			// 打印慢查询阈值
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	})

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

type gormLoggerFunc func(msg string, fields ...logger.Field)

// @func: Printf
// @date: 2023-11-20 23:20:10
// @brief: 函数衍生类型实现接口
// @author: Kewin Li
// @receiver g
// @param s
// @param i
func (g gormLoggerFunc) Printf(msg string, fields ...interface{}) {
	g(msg, logger.Field{"args", fields})
}
