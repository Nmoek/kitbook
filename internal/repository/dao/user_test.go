// Package dao
// @Description: 测试数据库相关接口
package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	msqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

// @func: TestGormUserDao_Insert
// @date: 2023-11-05 03:18:22
// @brief: 单元测试-数据库插入数据
// @author: Kewin Li
// @param t
func TestGormUserDao_Insert(t *testing.T) {
	nowMs := time.Now().UnixMilli()
	id := int64(1)
	testUser := User{
		Id:       id,
		Email:    sql.NullString{String: "123@qq.com", Valid: true},
		Phone:    sql.NullString{String: "1851561615", Valid: true},
		Password: "Ljk741610",
		Birthday: nowMs,
		AboutMe:  "这是一个mock测试",
		Ctime:    nowMs,
		Utime:    nowMs,
	}
	testCases := []struct {
		name string

		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入数据成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockRes)

				return db
			},

			ctx:  context.Background(),
			user: testUser,
		},
		{
			name: "唯一索引冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO .*").WillReturnError(&msqlDriver.MySQLError{Number: 1062})

				return db
			},

			ctx:  context.Background(),
			user: testUser,

			wantErr: ErrDuplicateUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true, //跳过版本查询
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true, //
			})

			assert.NoError(t, err)

			dao := NewGormUserDao(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
