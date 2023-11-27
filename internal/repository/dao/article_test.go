// Package dao
// @Description: 单元测试-帖子模块
package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

// @func: TestGormArticleDao_Sync
// @date: 2023-11-27 12:48:44
// @brief: 帖子发表-数据同步-dao层
// @author: Kewin Li
// @param t
func TestGormArticleDao_Sync(t *testing.T) {

	nowMs := time.Now().UnixMilli()
	testArt := Article{
		Id:       1,
		Title:    "发表的标题",
		Content:  "发表的内容",
		AuthorId: 123,
		Ctime:    nowMs,
		Utime:    nowMs,
	}

	testCases := []struct {
		name string

		mock func(t *testing.T) *sql.DB

		art Article

		wantErr error
		wantId  int64
	}{
		{
			name: "新建帖子, 两库保存成功, 同步成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				mock.ExpectBegin().WillReturnError(nil)

				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO *.").WillReturnResult(mockRes)

				return db
			},
			art:    testArt,
			wantId: 1,
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

			dao := NewGormArticleDao(db)
			id, err := dao.Sync(context.Background(), testArt)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)

		})
	}
}
