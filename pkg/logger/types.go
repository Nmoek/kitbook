// Package logger
// @Description: 日志模块
package logger

type Logger interface {
	DEBUG(msg string, args ...Field)
	INFO(msg string, args ...Field)
	WARN(msg string, args ...Field)
	ERROR(msg string, args ...Field)
}

type Field struct {
	Key string
	Val any
}
