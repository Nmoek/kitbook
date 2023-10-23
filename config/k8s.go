//go:build k8s

// Package config
// @Description: k8s的配置
package config

var Config = config{
	DB: DBconfig{
		DSN: "root:root@tcp(kitbook-mysql:3308)/kitbook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: Redisconfig{
		Addr: "kitbook-redis:6380",
	},
}
