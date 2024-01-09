// Package web
// @Description: 集成测试-帖子模块（使用测试套件）
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"kitbook/internal/domain"
	startup2 "kitbook/internal/integration/startup"
	"kitbook/internal/repository/dao"
	ijwt "kitbook/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ArticleMongoDBHandlerSuite
// @Description: article_mongoDB测试套件
type ArticleMongoDBHandlerSuite struct {
	suite.Suite
	server     *gin.Engine
	mdb        *mongo.Database
	produceCol *mongo.Collection
	liveCol    *mongo.Collection
}

// @func: SetupTest
// @date: 2023-11-22 23:43:00
// @brief: 测试套件运行前的准备
// @author: Kewin Li
// @receiver a
func (a *ArticleMongoDBHandlerSuite) SetupSuite() {

	a.mdb = startup2.InitMongoDB()
	err := dao.InitCollection(a.mdb)
	assert.NoError(a.T(), err)

	node, err := snowflake.NewNode(1)
	assert.NoError(a.T(), err)

	a.produceCol = a.mdb.Collection("articles")
	a.liveCol = a.mdb.Collection("published_articles")
	hdl := startup2.NewArticleHandler(dao.NewMongoDBArticleDAO(a.mdb, node))

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
func (a *ArticleMongoDBHandlerSuite) TearDownTest() {
	t := a.T()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 清空制作库
	_, err := a.produceCol.DeleteMany(ctx, bson.D{})
	assert.NoError(t, err)

	// 情况线上库
	_, err = a.liveCol.DeleteMany(ctx, bson.D{})
	assert.NoError(t, err)
}

// @func: TestEdit
// @date: 2023-12-02 21:36:40
// @brief: 编辑后保存，不发表
// @author: Kewin Li
// @receiver a
func (a *ArticleMongoDBHandlerSuite) TestEdit() {
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
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err := a.produceCol.FindOne(ctx, bson.D{{"author_id", int64(123)}}).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				art.Ctime = 0
				art.Utime = 0
				art.Id = 0
				assert.Equal(t, dao.Article{
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
			name: "修改帖子, 保存成功",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				_, err := a.produceCol.InsertOne(ctx, dao.Article{
					Id:       2,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				})

				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. 验证数据存入数据库
				// 2. 及时清理数据
				var art dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err := a.produceCol.FindOne(ctx, bson.D{{"id", 2}}).Decode(&art)
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
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				_, err := a.produceCol.InsertOne(ctx, &dao.Article{
					Id:       int64(3),
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					AuthorId: 456,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				})

				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. 验证数据存入数据库
				// 2. 及时清理数据
				var art dao.Article
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				err := a.produceCol.FindOne(ctx, bson.D{{"id", 3}}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:       int64(3),
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					AuthorId: 456,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}, art)

			},

			art: Article{
				Id:      int64(3),
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

			// 雪花算法生成的ID并无法准确知道
			if tc.wantRes.Data > 0 {
				assert.True(t, res.Data > 0)
				assert.Equal(t, tc.wantRes.Msg, res.Msg)
			}

			if tc.wantRes.Data == -1 {
				assert.Equal(t, tc.wantRes.Msg, res.Msg)
			}

		})
	}

}

// @func: TestPublish
// @date: 2023-12-02 19:53:58
// @brief: 帖子发表
// @author: Kewin Li
// @receiver a
func (a *ArticleMongoDBHandlerSuite) TestPublish() {
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
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				var artProduce dao.Article
				err := a.produceCol.FindOne(ctx, bson.M{
					"author_id": 123,
				}).Decode(&artProduce)

				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime > 0)
				assert.True(t, artProduce.Utime > 0)
				assert.True(t, artProduce.Id > 0)
				artProduce.Id = 0
				artProduce.Ctime = 0
				artProduce.Utime = 0
				assert.Equal(t, dao.Article{
					Title:    "新发表的标题",
					Content:  "新发表的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
				}, artProduce)

				var artLive dao.PublishedArticle
				err = a.produceCol.FindOne(ctx, bson.M{
					"author_id": 123,
				}).Decode(&artLive)
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime > 0)
				assert.True(t, artLive.Utime > 0)
				assert.True(t, artLive.Id > 0)
				artLive.Id = 0
				artLive.Ctime = 0
				artLive.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
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
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				art := dao.Article{
					Id:       2,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 123,
					Ctime:    666,
					Utime:    666,
				}

				_, err := a.produceCol.InsertOne(ctx, art)
				assert.NoError(t, err)

				_, err = a.liveCol.InsertOne(ctx, dao.PublishedArticle(art))
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				var artProduce dao.Article
				err := a.produceCol.FindOne(ctx, bson.M{"id": 2}).Decode(&artProduce)
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
				err = a.liveCol.FindOne(ctx, bson.M{"id": 2}).Decode(&artLive)
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

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				art := dao.Article{
					Id:       3,
					Title:    "修改前的标题",
					Content:  "修改前的内容",
					Status:   domain.ArticleStatusPublished,
					AuthorId: 456,
					Ctime:    777,
					Utime:    777,
				}

				_, err := a.produceCol.InsertOne(ctx, art)
				assert.NoError(t, err)

				_, err = a.liveCol.InsertOne(ctx, dao.PublishedArticle(art))
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				var artProduce dao.Article
				err := a.produceCol.FindOne(ctx, bson.M{"id": 3}).Decode(&artProduce)
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
				err = a.liveCol.FindOne(ctx, bson.M{"id": 3}).Decode(&artLive)
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
			if tc.wantRes.Data > 0 {

				assert.Equal(t, tc.wantRes.Msg, res.Msg)
				assert.True(t, tc.wantRes.Data > 0)
			}

			if tc.wantRes.Data == -1 {
				assert.Equal(t, tc.wantRes.Msg, res.Msg)
				assert.Equal(t, tc.wantRes.Data, res.Data)

			}

		})

	}

}

// @func: TestWithdraw
// @date: 2023-12-03 17:09:29
// @brief: 帖子撤回
// @author: Kewin Li
// @receiver a
func (a *ArticleMongoDBHandlerSuite) TestWithdraw() {
	t := a.T()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

				_, err := a.produceCol.InsertOne(ctx, art)
				assert.NoError(t, err)

				_, err = a.liveCol.InsertOne(ctx, dao.PublishedArticle(art))
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.produceCol.FindOne(ctx, bson.M{"id": 1}).Decode(&artProduce)
				assert.NoError(t, err)
				assert.True(t, artProduce.Ctime == 666)
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
				err = a.liveCol.FindOne(ctx, bson.M{"id": 1}).Decode(&artLive)
				assert.NoError(t, err)
				assert.True(t, artLive.Ctime == 666)
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

				_, err := a.produceCol.InsertOne(ctx, art)
				assert.NoError(t, err)

				_, err = a.liveCol.InsertOne(ctx, dao.PublishedArticle(art))

				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				var artProduce dao.Article
				err := a.produceCol.FindOne(ctx, bson.M{"id": 3}).Decode(&artProduce)
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
				err = a.liveCol.FindOne(ctx, bson.M{"id": 3}).Decode(&artLive)
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

			if tc.wantRes.Data > 0 {
				assert.Equal(t, tc.wantRes.Msg, res.Msg)
				assert.True(t, tc.wantRes.Data > 0)

			}

			if tc.wantRes.Data == -1 {
				assert.Equal(t, tc.wantRes.Msg, res.Msg)
				assert.Equal(t, tc.wantRes.Data, res.Data)

			}

		})

	}

}

// @func: TestArticleHandler
// @date: 2023-11-22 23:45:54
// @brief: 测试套件入口
// @author: Kewin Li
// @param t
func TestArticleMongoDBHandler(t *testing.T) {

	suite.Run(t, &ArticleMongoDBHandlerSuite{})

}
