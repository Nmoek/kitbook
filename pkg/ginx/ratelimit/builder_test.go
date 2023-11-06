// Package ratelimit
// @Description: 单元测试-限流器
package ratelimit

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareBuilder_limit(t *testing.T) {

	cmd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mdl := NewMiddlewareBuilder(cmd, 3*time.Second, 1).Build()

	server := gin.Default()
	server.Use(mdl)
	response := &TestResponseHandle{}
	server.GET("/limit", response.TestHandle)

	req, err := http.NewRequest(http.MethodGet, "/limit", nil)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()

	t.Run("不限流", func(t *testing.T) {

		server.ServeHTTP(recorder, req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, 1, response.cnt)
	})

	recorder = httptest.NewRecorder()
	t.Run("限流, 返回错误码", func(t *testing.T) {

		server.ServeHTTP(recorder, req)
		require.NoError(t, err)
		assert.Equal(t, 1, response.cnt)
		assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
	})

	recorder = httptest.NewRecorder()
	t.Run("等待窗口后, 不限流", func(t *testing.T) {
		time.Sleep(3 * time.Second) //等待窗口恢复

		server.ServeHTTP(recorder, req)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, 2, response.cnt)
	})

}

type TestResponseHandle struct {
	cnt int
}

func (t *TestResponseHandle) TestHandle(ctx *gin.Context) {
	t.cnt++
	ctx.String(http.StatusOK, "OK")
}
