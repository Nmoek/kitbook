// Package sms
// @Description: 短信模块
package sms

import "context"

type Service interface {
	//Send(ctx context.Context, param SendParam)
	Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error
}

// SendParam
// @Description: 短信发送入参数
//type SendParam struct {
//	AppId         string //一般是固定的
//	SignName      string //一般是固定的
//	SenderId      string
//	SessionCtx    string
//	ExtendCode    string
//	TemplateParam []string
//	TemplateId    string
//	PhoneNumber   []string
//}
