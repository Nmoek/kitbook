package startup

import (
	"kitbook/internal/service/oauth2/wechat"
)

func InitWechatService() wechat.Service {

	appID := "kitbook"
	appSecret := "123456"
	return wechat.NewService(appID, appSecret)
}
