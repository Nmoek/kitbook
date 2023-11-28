package logger

// 用户模块
const (
	LOG_USER_SIGNUP = iota
	LOG_USER_LOGIN
	LOG_USER_LOGINSMS
	LOG_USER_EDIT
	LOG_USER_PROFILE
	LOG_USER_REFRESHTOKEN
	LOG_USER_SENDCODE
	LOG_USER_LOGOUT
)

// 微信模块
const (
	LOG_WECHAT_AUTH2URL = iota
	LOG_WECHAT_CALLBACK
)

// 帖子模块
const (
	LOG_ART_EDIT = iota
	LOG_ART_PUBLISH
	LOG_ART_WITHDRAW
)

// 用户模块报错key
var UserLogMsgKey = map[int]string{
	LOG_USER_SIGNUP:       "user_signup_log",
	LOG_USER_LOGIN:        "user_login_log",
	LOG_USER_LOGINSMS:     "user_loginsms_log",
	LOG_USER_EDIT:         "user_edit_log",
	LOG_USER_PROFILE:      "user_profile_log",
	LOG_USER_REFRESHTOKEN: "user_refresh_log",
	LOG_USER_LOGOUT:       "user_logout_log",
}

// 微信模块报错key
var WechatLogMsgKey = map[int]string{
	LOG_WECHAT_AUTH2URL: "wechat_auth2url_log",
	LOG_WECHAT_CALLBACK: "wechat_callback_log",
}

// 帖子模块报错
var ArticleLogMsgKey = map[int]string{
	LOG_ART_EDIT:     "art_edit_log",
	LOG_ART_PUBLISH:  "art_publish_log",
	LOG_ART_WITHDRAW: "art_LOG_ART_WITHDRAW_log",
}
