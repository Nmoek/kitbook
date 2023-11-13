package ioc

import (
	"kitbook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {
	//appID := os.Getenv("")
	//appID, ok := os.LookupEnv("")
	//if !ok {
	//	panic("not find PATH ENV: WECHAT_APP_ID")
	//}

	//appSecret := os.LookupEnv("WECHAT_APP_SECRET")
	//appSecret := os.Getenv("")

	//if !ok {
	//	panic("not find PATH ENV: WECHAT_APP_SECRET")
	//}

	appID := "kitbook"
	appSecret := "123456"
	return wechat.NewService(appID, appSecret)
}
