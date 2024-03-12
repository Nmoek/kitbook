// Package ioc
// @Description: Redis初始化
package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRlockClient(client redis.Cmdable) *rlock.Client {
	return rlock.NewClient(client)
}

func InitRedis() redis.Cmdable {

	//return rdb.NewClient(&rdb.Options{
	//	Addr:     config.Config.Redis.Addr,
	//	Password: config.Config.Redis.Password, // no password docs
	//	DB:       0,                            // use default DB
	//})

	//return redis.NewClient(&redis.Options{
	//	Addr:     viper.GetString("redis.addr"),
	//	Password: viper.GetString("redis.password"),
	//})

	type Config struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	}

	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password docs
		DB:       0,            // use default DB
	})
}
