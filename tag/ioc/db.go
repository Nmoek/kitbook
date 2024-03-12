package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"kitbook/pkg/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
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
			//Colorful:      true,
			//SlowThreshold: 0, // 打印慢查询阈值
			//LogLevel:      glogger.Info,
		}),
	})

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
