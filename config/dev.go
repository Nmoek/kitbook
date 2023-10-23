//go:build !k8s

// Package config
// @Description: 本地配置
package config

var Config = config{
	DB: DBconfig{
		DSN: "root:root@tcp(127.0.0.1:13316)/kitbook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: Redisconfig{
		Addr: "localhost:6379",
	},
}
