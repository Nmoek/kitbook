// Package auth
// @Description: 短信服务使用的认证
package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"kitbook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

func NewSMSService(svc sms.Service, key []byte) *SMSService {
	return &SMSService{
		svc: svc,
		key: key,
	}
}

func (s *SMSService) Send(ctx context.Context, tplToken string, args []string, phoneNumber []string) error {
	var claims SMSClaims

	// 通过token解析出template
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}

	return s.svc.Send(ctx, claims.Tpl, args, phoneNumber)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
	// 额外的认真字段
}
