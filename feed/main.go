package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	// 初始化配置模块
	initViper()
	initPrometheus()
	// 初始化Web服务
	app := InitApp()

	app.rpcServer.Serve()

}

func initViper() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/dev.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func initPrometheus() {
	go func() {
		// 专门给prometheus用的端口
		http.Handle("metrics", promhttp.Handler())
		err := http.ListenAndServe("localhost:8082", nil)
		if err != nil {
			panic(err)
		}
	}()
}
