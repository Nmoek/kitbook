// Package ioc
// @Description: Redis初始化
package ioc

import (
	rdb "github.com/redis/go-redis/v9"
	"kitbook/config"
)

func InitRedis() rdb.Cmdable {

	return rdb.NewClient(&rdb.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password, // no password docs
		DB:       0,                            // use default DB
	})
}
