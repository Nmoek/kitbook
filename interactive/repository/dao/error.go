package dao

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrRepeatCancel   = errors.New("重复取消点赞/收藏")
)
