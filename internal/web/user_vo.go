package web

import (
	"kitbook/internal/domain"
	"time"
)

// ProfileVo
// @Description: 前端响应-用户信息
type ProfileVo struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}

func ConvertsProfileVo(user *domain.User) ProfileVo {
	return ProfileVo{
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Birthday: user.Birthday.Format(time.DateOnly),
		AboutMe:  user.AboutMe,
	}
}
