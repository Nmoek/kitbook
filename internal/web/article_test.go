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
)

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
				Data: float64(777), // json中数字转为go类型默认是float64
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
			hdl := NewArticleHandler(svc, logger.NewNopLogger())

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
