// Package smsTencent
// @Description: 腾讯云短信服务实现
package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	smsTencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
)

type Service struct {
	Client   *smsTencent.Client
	appId    string //固定不变
	signName string //固定不变
}

func NewServiceTencent(cli *smsTencent.Client, appId string, signature string) *Service {
	return &Service{
		Client:   cli,
		appId:    appId,
		signName: signature,
	}
}

// @func: Send
// @date: 2023-10-27 01:43:10
// @brief: 腾讯云-短信发送实现
// @author: Kewin Li
// @receiver s
// @param ctx
// @param templateId
// @param args
// @param phoneNumber
func (s *Service) Send(ctx context.Context, templateId string, args []string, phoneNumber []string) error {

	request := smsTencent.NewSendSmsRequest()
	//TODO: 链路数据，后续进行讲解
	request.SetContext(ctx)
	/* 基本类型的设置:
	 * SDK采用的是指针风格指定参数，即使对于基本类型你也需要用指针来对参数赋值。
	 * SDK提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台: https://console.cloud.tencent.com/smsv2
	 * smsTencent helper: https://cloud.tencent.com/document/product/382/3773 */
	/* 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666 */
	request.SmsSdkAppId = common.StringPtr(s.appId)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，签名信息可登录 [短信控制台] 查看 */
	request.SignName = common.StringPtr(s.signName)
	/* 国际/港澳台短信 SenderId: 中国大陆地区短信填空，默认未开通，如需开通请联系 [smsTencent helper] */
	//request.SenderId = common.StringPtr("")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	//request.SessionContext = common.StringPtr("")
	/* 短信码号扩展号: 默认未开通，如需开通请联系 [smsTencent helper] */
	//request.ExtendCode = common.StringPtr("")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(args)
	/* 模板 ID: 必须填写已审核通过的模板 ID。模板ID可登录 [短信控制台] 查看 */
	request.TemplateId = common.StringPtr(templateId)
	/* 下发手机号码，采用 E.164 标准，+[国家或地区码][手机号]
	 * 示例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	request.PhoneNumberSet = common.StringPtrs(phoneNumber)
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := s.Client.SendSms(request)

	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//TODO: 日志打印
		//fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。实际代码中可以加入其他的处理。
	if err != nil {
		//TODO: 日志打印
		//panic(err)
		fmt.Printf("err: %s \n", err)
		return err
	}

	var flag bool = false
	for _, st := range response.Response.SendStatusSet {
		if st != nil {
			status := *st
			if status.Code == nil || *status.Code != "ok" {
				//TODO: 有一个号码失败就全部返回是不太合理的，埋点打印，重发短信。
				fmt.Printf("腾讯云短信发送失败! [%s]: %s \n", *status.PhoneNumber, *status.Message)
				flag = true
			}

		}
	}

	if flag {
		return fmt.Errorf("message send err")
	}

	b, err := json.Marshal(response.Response)
	// 打印返回的json字符串
	fmt.Printf("sms Tencent: %s \n", string(b))

	return err
}
