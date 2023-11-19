// Package middlewares
// @Description: 利用Gin的middleware打印日志
package middlewares

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LogMiddlewareBuilder struct {
	logFunc       func(ctx context.Context, al AccessLog)
	pathMaxlen    int
	bodyMaxlen    int
	allowReqBody  bool
	allowRespBody bool
}

func NewLogMiddlewareBuilder(f func(ctx context.Context, al AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFunc: f,
	}
}

// 这里设计全为链式调用
func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

// 这里设计全为链式调用
func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// TODO: 这里着重考虑的是长度校验，防黑客
		if len(path) > l.pathMaxlen {
			path = path[:l.pathMaxlen]
		}

		method := ctx.Request.Method

		al := AccessLog{
			Path:   path,
			Method: method,
		}

		if l.allowReqBody {
			// 坑点, 这里的Body是一个stream对象，只能读取一次
			reqBody, _ := ctx.GetRawData()
			if len(reqBody) > l.bodyMaxlen {
				al.ReqBody = string(reqBody[:l.bodyMaxlen])
			} else {
				al.ReqBody = string(reqBody)
			}

			// 用完之后放回去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		start := time.Now()

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				al:             &al,
				ResponseWriter: ctx.Writer,
			}
		}

		// 防止显式调用Next时的崩溃
		defer func() {
			al.Duration = time.Since(start)

			l.logFunc(ctx, al)

		}()

		// 直接执行下一个middleware
		ctx.Next()
	}
}

type AccessLog struct {
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	ReqBody  string        `json:"req_body"`
	RespBody string        `json:"resp_body"`
	Status   int           `json:"status"`
	Duration time.Duration `json:"duration"`
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (r *responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.al.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// 该接口无用, 为接口实现做保证
func (r *responseWriter) WriteString(s string) (int, error) {
	return r.ResponseWriter.WriteString(s)
}
