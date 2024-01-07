package ioc

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"kitbook/pkg/logger"
)

func InitLogger() logger.Logger {
	//  zap日志配置加载
	// 开发环境使用development
	cfg := zap.NewDevelopmentConfig()
	// 生产环境使用release
	//cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true //关闭堆栈追踪
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.NewZapLogger(l)
}
