// Package web
// @Description: 用户模块-单元测试
package web

import (
	regexp "github.com/dlclark/regexp2"
	"testing"
)

// @func: TestPhoneRegExp
// @date: 2023-10-29 21:33:08
// @brief: 测试手机号校验
// @author: Kewin Li
// @param t
func TestPhoneRegExp(t *testing.T) {

	const phoneRegexPattern = "(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}"
	phoneRegExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	ok, err := phoneRegExp.MatchString("15662850585")

	t.Logf("ok:%v, err:%v \n", ok, err)
}
