// Package service
// @Description: 帖子服务-单元测试
package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	repomocks "kitbook/internal/repository/mocks"
	"kitbook/pkg/logger"
	"testing"
)

// @func: TestNormalArticleService_Publish
// @date: 2023-11-26 01:25:33
// @brief: 单元测试-帖子发表
// @author: Kewin Li
// @param t
func TestNormalArticleService_Publish(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (
			repository.ArticleAuthorRepository,
			repository.ArticleReaderRepository)

		art domain.Article

		wantArtId int64
		wantErr   error
	}{
		{
			name: "新建帖子，并发表成功",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {

				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				return authorRepo, readerRepo
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantArtId: int64(1),
		},
		{
			name: "新建帖子，并发表失败",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {

				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("发表失败"))

				return authorRepo, readerRepo
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr:   errors.New("发表失败"),
			wantArtId: int64(1),
		},
		{
			name: "新建帖子，并保存失败",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {

				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), errors.New("制作库保存失败"))

				return authorRepo, readerRepo
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr: errors.New("制作库保存失败"),
		},
		{
			name: "修改帖子，并首次发表成功",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {

				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				return authorRepo, readerRepo
			},

			art: domain.Article{
				Id:      2,
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantArtId: int64(2),
		},
		{
			name: "修改帖子，并保存失败",
			mock: func(ctrl *gomock.Controller) (
				repository.ArticleAuthorRepository,
				repository.ArticleReaderRepository) {

				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      3,
					Title:   "发表的标题",
					Content: "发表的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("制作库保存失败"))

				return authorRepo, readerRepo
			},

			art: domain.Article{
				Id:      3,
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr:   errors.New("制作库保存失败"),
			wantArtId: int64(3),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewNormalArticleServiceV1(authorRepo, readerRepo, logger.NewNopLogger())

			artId, err := svc.Publish(context.Background(), tc.art)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArtId, artId)
		})
	}
}
