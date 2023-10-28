// Package repository
// @Description: 验证码功能
package repository

import (
	"context"
	"kitbook/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

// @func: Set
// @date: 2023-10-28 21:53:48
// @brief: 转发模块-验证码发送
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @param code
// @return error
func (c *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

// @func: Verify
// @date: 2023-10-29 00:14:53
// @brief: 转发模块-验证码验证
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @param inputCode
// @return bool
// @return error
func (c *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)
}
