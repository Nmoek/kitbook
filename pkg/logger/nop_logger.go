// Package logger
// @Description: 不进行任何日志打印
package logger

type NopLogger struct {
}

func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

func (n *NopLogger) DEBUG(msg string, args ...Field) {
	//Nothing to do
}

func (n *NopLogger) INFO(msg string, args ...Field) {
	//Nothing to do

}

func (n *NopLogger) WARN(msg string, args ...Field) {
	//Nothing to do

}

func (n *NopLogger) ERROR(msg string, args ...Field) {
	//Nothing to do

}
