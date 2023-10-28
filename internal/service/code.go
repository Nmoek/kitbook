// Package service
// @Description: 验证码功能
package service

import (
	"context"
	"fmt"
	"kitbook/internal/repository"
	"kitbook/internal/repository/cache"
	"kitbook/internal/service/sms"
	"math/rand"
)

// 验证码模板, 根据想生成的数字模板随便填写, 形式风格会与你填写的保持一致
const templateId = "123456"

var ErrCodeSendTooMany = cache.ErrCodeSendTooMany
var ErrCodeVerifyCntTooMany = cache.ErrCodeVerifyCntTooMany
var ErrCodeNotRight = cache.ErrCodeNotRight

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
}

// @func: Send
// @date: 2023-10-28 20:43:54
// @brief: 随机生成验证码并发送
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @return error
func (c *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := c.generate()
	err := c.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	return c.sms.Send(ctx, templateId, []string{code}, []string{phone})
}

// @func: Verify
// @date: 2023-10-28 20:45:08
// @brief: 验证用户填写的验证码
// @author: Kewin Li
// @receiver c
// @param ctx
// @param biz
// @param phone
// @param inputCode
// @return bool
// @return error
func (c *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	if err == ErrCodeVerifyCntTooMany {
		//TODO: 验证码超过验证次数，日志埋点
		return false, nil
	}

	return ok, err
}

// @func: generate
// @date: 2023-10-28 21:58:24
// @brief: 生成6位随机验证码
// @author: Kewin Li
// @receiver c
// @return string
func (c *CodeService) generate() string {

	code := rand.Intn(10000000)
	return fmt.Sprintf("%06d", code)
}
