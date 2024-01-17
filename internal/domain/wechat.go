package domain

type WechatInfo struct {
	//	用户统一标识。针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的
	Unionid string
	//	授权用户唯一标识
	Openid string
}
