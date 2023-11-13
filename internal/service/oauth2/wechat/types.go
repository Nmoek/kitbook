// Package wechat
// @Description: 微信验证服务
package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"kitbook/internal/domain"
	"net/http"
	"net/url"
)

const autuURLPattern = `https: //open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
const accessTokenURLPattern = `https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`

var redirectURL = url.PathEscape(`https://meoying.com/oauth2/wechat/callback`)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechtInfo, error)
}

type service struct {
	appID     string
	appSecret string
	client    *http.Client
}

func NewService(appID string, appSecret string) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	// 第一个参数 appid
	// 第二个参数 redirect_uri
	// 第三个参数 state
	return fmt.Sprintf(autuURLPattern, s.appID, redirectURL, state), nil
}

// @func: VerifyCode
// @date: 2023-11-12 00:12:37
// @brief: 通过code获取access_token
// @author: Kewin Li
// @receiver s
// @param ctx
// @param code
// @return error
func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechtInfo, error) {
	// 第一个参数 appid
	// 第二个参数 secret
	// 第三个参数 code
	accessTokenUrl := fmt.Sprintf(accessTokenURLPattern, s.appID, s.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechtInfo{}, err
	}

	response, err := s.client.Do(req)
	if err != nil {
		return domain.WechtInfo{}, err
	}

	var res Result
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		return domain.WechtInfo{}, err
	}

	return domain.WechtInfo{
		Unionid: res.Unionid,
		Openid:  res.Openid,
	}, nil
}

type Result struct {
	//	接口调用凭证
	AccessToken string `json:"access_token"`

	//access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`

	//	用户刷新access_token
	RefreshToken string `json:"refresh_token"`

	//	授权用户唯一标识
	Openid string `json:"openid"`

	//用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`

	//	用户统一标识。针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的
	Unionid string `json:"unionid"`

	// 错误码
	ErrCode string `json:"errcode"`
	// 错误信息
	ErrMsg string `json:"errmsg"`
}
