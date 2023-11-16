// Package jwt
// @Description: 基于redis的jwt交互
package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type RedisJWTHandler struct {
	cmd           redis.Cmdable
	signingMethod jwt.SigningMethod
	/// ssid过期时间
	expiration time.Duration
}

func NewRedisJWTHandler(cmd redis.Cmdable) JWTHandler {
	return &RedisJWTHandler{
		cmd:           cmd,
		expiration:    time.Hour * 24 * 7,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// @func: SetJWTToken
// @date: 2023-11-15 02:15:33
// @brief: 设置token
// @author: Kewin Li
// @receiver r
// @param ctx
// @param id
// @param ssid
// @return error
func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, id int64, ssid string) error {

	// TODO: 一种方案在这里刷新长token的生效时间
	//err := j.setRefreshToken(ctx, id)
	//if err != nil {
	//	ctx.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}

	// 设置JWT
	token := jwt.NewWithClaims(r.signingMethod, UserClaims{
		UserID:    id,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 30min过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	})

	tokenStr, err := token.SignedString([]byte(TokenPrivateKey))
	if err != nil {

		return err
	}

	ctx.Header("x-jwt-token", tokenStr)

	return nil
}

// @func: setTokenWithSsid
// @date: 2023-11-15 02:17:27
// @brief: 设置带ssid的长短token(登出使用)
// @author: Kewin Li
// @receiver j
// @param ctx
// @param id
// @return error
func (r *RedisJWTHandler) SetTokenWithSsid(ctx *gin.Context, id int64) error {
	ssid := uuid.New().String()
	err := r.setRefreshToken(ctx, id, ssid)
	if err != nil {
		return err
	}

	err = r.SetJWTToken(ctx, id, ssid)
	if err != nil {
		return err
	}

	return nil
}

// @func: ClearToken
// @date: 2023-11-15 02:18:23
// @brief: 登出时清除ssid(设置无用token)
// @author: Kewin Li
// @receiver j
// @param ctx
// @return error
func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("user_token").(UserClaims)

	return r.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", r.expiration).Err()
}

// @func: CheckSsid
// @date: 2023-11-15 02:21:32
// @brief: 校验登出时的Ssid
// @author: Kewin Li
// @receiver r
// @param ctx
// @param ssid
// @return error
func (r *RedisJWTHandler) CheckSsid(ctx *gin.Context, ssid string) error {
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil {
		return err
	}

	if cnt > 0 {
		return errors.New("token已失效")
	}

	return nil
}

// @func: ExtractToken
// @date: 2023-11-15 02:12:16
// @brief: 从auth字段里拿回token进行解析
// @author: Kewin Li
// @receiver r
// @param ctx
// @return string
func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	// Bearer token
	authData := ctx.GetHeader("Authorization")
	if authData == "" {
		// Authorization字段未设置
		return authData
	}

	segs := strings.Split(authData, " ")
	if len(segs) < 2 {
		// token格式有误
		return ""
	}

	return segs[1]
}

// @func: setRefreshToken
// @date: 2023-11-13 23:48:55
// @brief: 设置长token
// @author: Kewin Li
// @receiver r
// @param ctx
// @param id
func (r *RedisJWTHandler) setRefreshToken(ctx *gin.Context, id int64, ssid string) error {

	token := jwt.NewWithClaims(r.signingMethod, RefreshClaims{
		UserID: id,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 七天过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.expiration)),
		},
	})

	tokenStr, err := token.SignedString([]byte(TokenRefreshKey))
	if err != nil {
		return err
	}

	ctx.Header("x-refresh-token", tokenStr)

	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID int64
	Ssid   string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
	Ssid      string
}
