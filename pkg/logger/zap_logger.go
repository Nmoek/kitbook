// Package logger
// @Description: 基于zap的日志封装
package logger

import "go.uber.org/zap"

// 与zap库对齐的报错码级别
const (
	DebugLevel  = -1
	InfoLevel   = 0
	WarnLevel   = 1
	ErrorLevel  = 2
	DPanicLevel = 3
	PanicLevel  = 4
	FatalLevel  = 5
)

const (
	LOG_SIGNUP = iota
	LOG_LOGIN
	LOG_LOGINSMS
	LOG_EDIT
	LOG_PROFILE
	LOG_SEND_SMSCODE
	LOG_LOGOUT
)

// 用户模块报错key
var UserLogMsgKey = map[int]string{
	LOG_SIGNUP:   "signup_log",
	LOG_LOGIN:    "login_log",
	LOG_LOGINSMS: "loginsms_log",
	LOG_EDIT:     "edit_log",
	LOG_PROFILE:  "profile_log",
	LOG_LOGOUT:   "logout_log",
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{
		logger: logger,
	}
}

func (z *ZapLogger) DEBUG(msg string, args ...Field) {
	z.logger.Debug(msg, z.toArgs(args)...)
}

func (z *ZapLogger) INFO(msg string, args ...Field) {
	z.logger.Info(msg, z.toArgs(args)...)
}

func (z *ZapLogger) WARN(msg string, args ...Field) {
	z.logger.Warn(msg, z.toArgs(args)...)
}

func (z *ZapLogger) ERROR(msg string, args ...Field) {
	z.logger.Error(msg, z.toArgs(args)...)
}

// @func: toArgs
// @date: 2023-11-18 17:38:49
// @brief: Field --> zap.Field
// @author: Kewin Li
// @receiver z
// @param args
// @return []zap.Field
func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))

	for _, f := range args {
		res = append(res, zap.Any(f.Key, f.Val))
	}

	return res
}
