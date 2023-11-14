// Package jwt
// @Description: 将jwt相关操作抽象整合
package jwt

import "github.com/gin-gonic/gin"

const (
	TokenPrivateKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C5"
	TokenRefreshKey = "kAEpRBDAb1PlhOHdpHYelwdNIsjmJ5C2"
)

type JWTHandler interface {
	SetJWTToken(ctx *gin.Context, id int64, ssid string) error
	SetTokenWithSsid(ctx *gin.Context, id int64) error
	ClearToken(ctx *gin.Context) error
	CheckSsid(ctx *gin.Context, ssid string) error
	ExtractToken(ctx *gin.Context) string
}
