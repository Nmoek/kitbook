package ioc

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"kitbook/pkg/logger"
	"time"
)

func InitLogger() logger.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "/var/log/user.log", // 日志文件路劲
		MaxSize:    50,                  // 每个日志文件大小, 单位MB
		MaxBackups: 3,                   // 保留旧日志的最大个数
		MaxAge:     28,                  //保留旧日志的最大天数
		Compress:   true,                //是否压缩旧日志文件
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(lumberjackLogger),
		zapcore.DebugLevel, // 设置日志级别
	)

	//  zap日志配置加载
	// 开发环境使用development
	//cfg := zap.NewDevelopmentConfig()
	//// 生产环境使用release
	////cfg := zap.NewProductionConfig()
	//cfg.DisableStacktrace = true //关闭堆栈追踪
	//err := viper.UnmarshalKey("log", &cfg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//l, err := cfg.Build()
	//if err != nil {
	//	panic(err)
	//}

	// 此处演示ELK
	l := zap.New(core, zap.AddCaller())
	res := logger.NewZapLogger(l)
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		for t := range ticker.C {
			res.INFO("模拟输出日志", logger.Field{"time", t.String()})
		}
	}()
	return res
}
