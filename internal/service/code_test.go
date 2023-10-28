// Package service
// @Description: 验证码功能-单元测试
package service

import (
	"fmt"
	"math/rand"
	"testing"
)

// @func: TestCodeGenerate
// @date: 2023-10-28 22:01:23
// @brief: 单元测试-生成随机验证码
// @author: Kewin Li
// @param t
func TestCodeGenerate(t *testing.T) {
	t.Log(fmt.Sprintf("%06d", rand.Intn(1000000)))
}
