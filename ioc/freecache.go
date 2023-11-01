// Package ioc
// @Description: freecache本地缓存
package ioc

import "github.com/coocood/freecache"

const localCacheSize = 200 * 1024 * 1024

func InitFreeCache() *freecache.Cache {
	return freecache.NewCache(localCacheSize)
}
