// Package web
// @Description: web层-帖子模块-单元测试
package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/domain"
	"kitbook/internal/service"
	svcmocks "kitbook/internal/service/mocks"
	ijwt "kitbook/internal/web/jwt"
	"kitbook/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// @func: TestArticleHandler_Publish
// @date: 2023-11-29 01:03:05
// @brief: 帖子发表-单元测试
// @author: Kewin Li
// @param t
func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleService

		reqBody  string
		wantCode int
		wantRes  Result
	}{
		{
			name: "新建帖子，并发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "发表标题",
					Content: "发表内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				return svc
			},

			reqBody: `{
"title": "发表标题",
"content": "发表内容"
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg:  "发表成功",
				Data: float64(1), // json中数字转为go类型默认是float64
			},
		},
		{
			name: "已有帖子，并发表成功",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      666,
					Title:   "发表标题",
					Content: "发表内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(666), nil)

				return svc
			},

			reqBody: `{
"id": 666,
"title": "发表标题",
"content": "发表内容"
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg:  "发表成功",
				Data: float64(666), // json中数字转为go类型默认是float64
			},
		},
		{
			name: "已有帖子，并发表失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      777,
					Title:   "发表标题",
					Content: "发表内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(777), errors.New("未知错误"))

				return svc
			},

			reqBody: `{
"id": 777,
"title": "发表标题",
"content": "发表内容"
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg:  "发表失败",
				Data: float64(-1), // json中数字转为go类型默认是float64
			},
		},
		{
			name: "Bind错误",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				return svc
			},

			reqBody: `{
"title": "发表标题",
"content": "发表内容"asdasdsads
}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tc.mock(ctrl)
			hdl := NewArticleHandler(svc, nil, logger.NewNopLogger())

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user_token", ijwt.UserClaims{
					UserID: 123,
				})
			})
			hdl.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewReader([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.wantCode, recorder.Code)
			if recorder.Code != http.StatusOK {
				return
			}

			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

// @func: TestArticleHandler_Withdraw
// @date: 2023-11-29 01:03:39
// @brief: 帖子撤回-单元测试
// @author: Kewin Li
// @param t
func TestArticleHandler_Withdraw(t *testing.T) {
	// TODO: 撤回帖子单元测试
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleService
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

// @func: TestArticleHandler_List
// @date: 2023-12-04 00:58:47
// @brief: 创作者列表
// @author: Kewin Li
// @param t
func TestArticleHandler_List(t *testing.T) {

	now := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.ArticleService

		reqBody string

		wantCode int
		wantRes  Result
	}{
		{
			name: "创作者列表查询成功",

			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)

				svc.EXPECT().GetByAuthor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]domain.Article{
						{
							Id:     1,
							Author: domain.Author{Id: 123},
							Status: domain.ArticleStatusPublished,
							Ctime:  now,
							Utime:  now,
						}, {
							Id:     2,
							Author: domain.Author{Id: 123},
							Status: domain.ArticleStatusPublished,
							Ctime:  now,
							Utime:  now,
						}, {
							Id:     3,
							Author: domain.Author{Id: 123},
							Status: domain.ArticleStatusPublished,
							Ctime:  now,
							Utime:  now,
						},
					}, nil)
				return svc
			},

			reqBody: `{
"limit": 3,
"offset": 0
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "查询成功",
				Data: []ArticleVo{
					{
						Id:     1,
						Status: domain.ArticleStatusPublished,
						Ctime:  now.Format(time.DateTime),
						Utime:  now.Format(time.DateTime),
					}, {
						Id:     2,
						Status: domain.ArticleStatusPublished,
						Ctime:  now.Format(time.DateTime),
						Utime:  now.Format(time.DateTime),
					}, {
						Id:     3,
						Status: domain.ArticleStatusPublished,
						Ctime:  now.Format(time.DateTime),
						Utime:  now.Format(time.DateTime),
					},
				},
			},
		},
		{
			name: "非创作者查询列表, 查询失败",

			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)

				svc.EXPECT().GetByAuthor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, service.ErrInvalidUpdate)
				return svc
			},

			reqBody: `{
"limit": 3,
"offset": 2
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "非法操作",
			},
		},
		{
			name: "系统错误, 查询失败",

			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)

				svc.EXPECT().GetByAuthor(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("其他错误"))
				return svc
			},

			reqBody: `{
"limit": 3,
"offset": 5
}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := tc.mock(ctrl)

			hdl := NewArticleHandler(svc, nil, logger.NewNopLogger())
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user_token", ijwt.UserClaims{
					UserID: 123,
				})
			})

			hdl.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/articles/list", bytes.NewReader([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			var res Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantRes.Data != nil && res.Data != nil {
				isResEqual(t, tc.wantRes.Data.([]ArticleVo), res.Data.([]any))

			}
		})
	}
}

func isResEqual(t *testing.T, arts []ArticleVo, datas []any) {

	for i, art := range arts {
		assert.Equal(t, float64(art.Id), datas[i].(map[string]any)["id"])
		if v, b := datas[i].(map[string]any)["author_id"]; b {
			assert.Equal(t, float64(art.AuthorId), v)
		}
		assert.Equal(t, float64(art.Status), datas[i].(map[string]any)["status"])
		assert.Equal(t, art.Ctime, datas[i].(map[string]any)["ctime"])
		assert.Equal(t, art.Utime, datas[i].(map[string]any)["utime"])
	}

}
