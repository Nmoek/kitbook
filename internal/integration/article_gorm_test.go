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
	"kitbook/internal/domain"
	startup2 "kitbook/internal/integration/startup"
	"kitbook/internal/repository/dao"
	ijwt "kitbook/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ArticleMongoDBHandlerSuite
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

	a.db = startup2.InitDB()
	hdl := startup2.NewArticleHandler(dao.NewGormArticleDao(a.db))

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
	a.db.Exec("truncate table `published_articles`")
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
		art Article

		wantCode int
		wantRes  Result[int64]
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
					Status:   domain.ArticleStatusUnpublished,
				}, art)
			},

			art: Article{
				Title:   "第一个帖子",
				Content: "第一个帖子的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "保存成功",
				Data: int64(1),
			},
		},
		{
			name: "修改未发表帖子, 保存成功",
			before: func(t *testing.T) {
				err := a.db.Create(&dao.Article{
					Id:       2,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished,
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
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    456,
				}, art)

			},

			art: Article{
				Id:      2,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
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
					Status:   domain.ArticleStatusPublished,
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
					Status:   domain.ArticleStatusPublished,
					AuthorId: 666,
					Ctime:    456,
				}, art)

			},

			art: Article{
				Id:      3,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
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
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

// @func: TestPublish
// @date: 2023-12-02 19:53:58
// @brief: 帖子发表
// @author: Kewin Li
// @receiver a
func (a *ArticleHandlerSuite) TestPublish() {
	t := a.T()

	testCases := []struct {
		name string

		before func(t *testing.T)

		after func(t *testing.T)

		art Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "新建帖子, 发表成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.db.Where("id = ?", 1).First(&artProduce).Error
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime > 0)
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "新发表的标题",
					Content:  "新发表的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.db.Where("id = ?", 1).First(&artLive).Error
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime > 0)
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					Id:       1,
					Title:    "新发表的标题",
					Content:  "新发表的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
				}, artLive)

			},

			art: Article{
				Title:   "新发表的标题",
				Content: "新发表的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "发表成功",
				Data: int64(1),
			},
		},
		{
			name: "修改帖子, 保存成功, 发表成功",
			before: func(t *testing.T) {

				art := dao.Article{
					Id:       2,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
					Ctime:    666,
					Utime:    666,
				}

				err := a.db.Create(art).Error
				assert.NoError(t, err)

				err = a.db.Create(dao.PublishedArticle(art)).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.db.Where("id = ?", 2).First(&artProduce).Error
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime > 666)
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "修改后的标题",
					Content:  "修改后的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.db.Where("id = ?", 2).First(&artLive).Error
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime > 666)
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					Id:       2,
					Title:    "修改后的标题",
					Content:  "修改后的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
				}, artLive)

			},

			art: Article{
				Id:      2,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "发表成功",
				Data: int64(2),
			},
		},
		{
			name: "修改别人帖子, 保存失败, 发表失败",
			before: func(t *testing.T) {

				err := a.db.Create(dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
					Ctime:    777,
					Utime:    777,
				}).Error
				assert.NoError(t, err)

				err = a.db.Create(dao.PublishedArticle{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
					Ctime:    777,
					Utime:    777,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.db.Where("id = ?", 3).First(&artProduce).Error
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime == 777)
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.db.Where("id = ?", 3).First(&artLive).Error
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime == 777)
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
				}, artLive)

			},

			art: Article{
				Id:      3,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "非法操作",
				Data: int64(-1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(&tc.art)
			assert.NoError(t, err)

			// 发表接口测试
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			// 执行请求
			a.server.ServeHTTP(recorder, req)

			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})

	}

}

// @func: TestWithdraw
// @date: 2023-12-02 23:21:55
// @brief: 帖子撤回
// @author: Kewin Li
// @receiver a
func (a *ArticleHandlerSuite) TestWithdraw() {
	t := a.T()

	testCases := []struct {
		name string

		before func(t *testing.T)

		after func(t *testing.T)

		art Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "撤回帖子, 撤回成功",
			before: func(t *testing.T) {

				art := dao.Article{
					Id:       1,
					Title:    "撤回帖子的标题",
					Content:  "撤回帖子的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
					Ctime:    666,
					Utime:    666,
				}

				err := a.db.Create(art).Error
				assert.NoError(t, err)

				err = a.db.Create(dao.PublishedArticle(art)).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.db.Where("id = ?", 1).First(&artProduce).Error
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime > 666)
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "撤回帖子的标题",
					Content:  "撤回帖子的内容",
					Status:   domain.ArticleStatusPrivate,
					AuthorId: 123,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.db.Where("id = ?", 1).First(&artLive).Error
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime > 666)
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					Id:       1,
					Title:    "撤回帖子的标题",
					Content:  "撤回帖子的内容",
					Status:   domain.ArticleStatusPrivate,
					AuthorId: 123,
				}, artLive)

			},

			art: Article{
				Id: 1,
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "撤回成功",
				Data: int64(1),
			},
		},
		{
			name: "撤回别人帖子, 撤回失败",
			before: func(t *testing.T) {

				art := dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
					Ctime:    777,
					Utime:    777,
				}

				err := a.db.Create(art).Error
				assert.NoError(t, err)

				err = a.db.Create(dao.PublishedArticle(art)).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.db.Where("id = ?", 3).First(&artProduce).Error
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime == 777)
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.db.Where("id = ?", 3).First(&artLive).Error
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime == 777)
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
				}, artLive)

			},

			art: Article{
				Id:      3,
				Title:   "修改后的标题",
				Content: "修改后的内容",
			},

			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "非法操作",
				Data: int64(-1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(&tc.art)
			assert.NoError(t, err)

			// 发表接口测试
			req, err := http.NewRequest(http.MethodPost, "/articles/withdraw", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			// 执行请求
			a.server.ServeHTTP(recorder, req)

			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantRes, res)

		})

	}

}

// @func: TestArticleHandler
// @date: 2023-11-22 23:45:54
// @brief: 测试套件入口
// @author: Kewin Li
// @param t
func TestArticleHandler(t *testing.T) {

	suite.Run(t, &ArticleHandlerSuite{})

}

type Article struct {
	Id      int64
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Result[T any]
// @Description: 使用泛型约束需要返回的内容的类型
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
