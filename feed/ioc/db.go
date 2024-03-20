// Package ioc
// @Description: 数据库初始化
package ioc

import (
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"kitbook/feed/repository/dao"
	"kitbook/pkg/gormx"
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
			//Colorful:      true,
			//SlowThreshold: 0, // 打印慢查询阈值
			//LogLevel:      glogger.Info,
		}),
	})

	if err != nil {
		panic(err)
	}

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "kitbook",
		RefreshInterval: 15, // 多久获取一次连接池的状态
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))

	if err != nil {
		panic(err)
	}

	/*prometheus 初始化*/
	cb := gormx.NewCallbacks(prometheus2.SummaryOpts{
		Namespace: "kewin",
		Subsystem: "kitbook",
		Name:      "gorm_db",
		Help:      "统计gorm查询库的执行时间",
		ConstLabels: map[string]string{
			"instance_id": "my_instance",
		},

		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})

	err = db.Use(cb)
	if err != nil {
		panic(err)
	}

	/*OPTL接入gorm*/
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
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
