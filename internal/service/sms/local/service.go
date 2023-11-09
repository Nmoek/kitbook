// Package local
// @Description: 本地发送验证码(验证)
package local

import (
	"context"
	"fmt"
	"kitbook/pkg/limiter"
)

type Service struct {
}

func (s *Service) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {

	fmt.Printf("验证码: %v \n", args)
	return nil
}

func NewService(limiter limiter.Limiter) *Service {
	return &Service{}
}
