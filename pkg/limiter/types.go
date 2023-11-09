// Package limiter
// @Description: 通用限流器插件
package limiter

import "context"

type Limiter interface {
	//  @interface_func: Limit
	//  @brief: 是否触发限流
	Limit(ctx context.Context, key string) (bool, error)
}

//type AccessLimiter interface {
//	Limiter
//}
//
//type SMSLimiter interface {
//	Limiter
//}
