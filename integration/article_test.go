// Package web
// @Description: 集成测试-帖子模块（使用测试套件）
package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"kitbook/integration/startup"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
	ijwt "kitbook/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ArticleHandlerSuite
// @Description: article测试套件
type ArticleHandlerSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

// @func: SetupTest
// @date: 2023-11-22 23:43:00
// @brief: 测试套件运行前的准备
// @author: Kewin Li
// @receiver a
func (a *ArticleHandlerSuite) SetupSuite() {

	a.db = startup.InitDB()
	hdl := startup.NewArticleHandler()
	//a.server = startup.InitWebServer()
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user_token", ijwt.UserClaims{
			UserID: 123,
		})
	})

	hdl.RegisterRoutes(server)
	a.server = server
}

// @func: TearDownSubTest
// @date: 2023-11-22 23:39:53
// @brief: 全部测试结束后执行的一个回调
// @author: Kewin Li
// @receiver a
func (a *ArticleHandlerSuite) TearDownTest() {
	// 注意: 不使用delete, 而是使用truncate完全清空表数据但不影响表结构
	a.db.Exec("truncate table `articles`")
}

// @func: TestArticleHandler_Edit
// @date: 2023-11-24 18:44:37
// @brief: 编辑后保存，不发表
// @author: Kewin Li
// @receiver a
func (a *ArticleHandlerSuite) TestEdit() {
	t := a.T()

	// TODO：为什么第三个用例无法成功执行？
	testCases := []struct {
		name string

		// 要提前准备的来自其他中间件的数据
		before func(t *testing.T)

		// 验证数据并及时清理数据
		after func(t *testing.T)

		// 前端传过来的数据
		art article

		wantCode int
		wantRes  result[int64]
	}{
		{
			name:   "新建帖子, 保存成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 1. 验证数据存入数据库
				// 2. 及时清理数据
				var art dao.Article
				err := a.db.Where("id = ?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "第一个帖子",
					Content:  "第一个帖子的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpblished,
				}, art)
			},

			art: article{
				Title:   "第一个帖子",
				Content: "第一个帖子的内容",
			},

			wantCode: http.StatusOK,
			wantRes: result[int64]{
				Msg:  "保存成功",
				Data: int64(1),
			},
		},
		{
			name: "修改已发表帖子, 保存成功",
			before: func(t *testing.T) {
				err := a.db.Create(&dao.Article{
					Id:       2,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPblished,
					Ctime:    456,
					Utime:    789,
				}).Error

				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. 验证数据存入数据库
				// 2. 及时清理数据
				var art dao.Article
				err := a.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 789)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "修改后的标题",
					Content:  "修改后的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpblished,
					Ctime:    456,
				}, art)

			},

			art: article{
				Id:      2,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: result[int64]{
				Msg:  "保存成功",
				Data: int64(2),
			},
		},
		{
			name: "修改别人帖子, 保存失败",
			before: func(t *testing.T) {
				err := a.db.Create(&dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPblished,
					AuthorId: 666,
					Ctime:    456,
					Utime:    789,
				}).Error

				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. 验证数据存入数据库
				// 2. 及时清理数据
				var art dao.Article
				err := a.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime == 789)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPblished,
					AuthorId: 666,
					Ctime:    456,
				}, art)

			},

			art: article{
				Id:      3,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: result[int64]{
				Msg:  "非法操作",
				Data: int64(-1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)

			// 准备http请求，接收http响应的recorder
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			// 执行，发出请求
			a.server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}

			// 通过json解码器解析出响应结果
			var res result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantRes, res)
		})
	}

	time.Sleep(2 * time.Second)
}

// @func: TestArticleHandler
// @date: 2023-11-22 23:45:54
// @brief: 测试套件入口
// @author: Kewin Li
// @param t
func TestArticleHandler(t *testing.T) {

	suite.Run(t, &ArticleHandlerSuite{})

}

type article struct {
	Id      int64
	Title   string `json:"title"`
	Content string `json:"content"`
}

// result[T any]
// @Description: 使用泛型约束需要返回的内容的类型
type result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
