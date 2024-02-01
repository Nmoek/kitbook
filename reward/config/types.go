// Package config
// @Description: 配置参数定义
package config

type config struct {
	DB    DBconfig
	Redis Redisconfig
}

type DBconfig struct {
	DSN string
}

type Redisconfig struct {
	Addr     string
	Password string
}
