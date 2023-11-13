package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const TokenPrivateKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C5"

type jwtHandler struct {
}

// @func: setJWTToken
// @date: 2023-10-29 23:06:05
// @brief: 设置JWT
// @author: Kewin Li
// @receiver h
// @param ctx
// @param id
func (j *jwtHandler) setJWTToken(ctx *gin.Context, id int64) {
	// 设置JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:    id,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	})

	tokenStr, err := token.SignedString([]byte(TokenPrivateKey))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误!")
		return
	}

	ctx.Header("x-jwt-token", tokenStr)
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}
