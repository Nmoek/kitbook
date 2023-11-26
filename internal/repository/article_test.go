// Package repository
// @Description: 数据转发层-帖子模块-单元测试
package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
	daomocks "kitbook/internal/repository/dao/mocks"
	"kitbook/pkg/logger"
	"testing"
)

// @func: TestCacheArticleRepository_Sync
// @date: 2023-11-26 19:15:02
// @brief: 单元测试-发表帖子
// @author: Kewin Li
// @param t
func TestCacheArticleRepository_Sync(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao)

		art     domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "新建帖子, 两库成功, 同步成功",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Create(gomock.Any(), dao.Article{
					Title:    "发表的标题",
					Content:  "发表的内容",
					AuthorId: 123,
				}).Return(int64(1), nil)

				readerDao.EXPECT().Upsert(gomock.Any(), dao.Article{
					Id:       1,
					Title:    "发表的标题",
					Content:  "发表的内容",
					AuthorId: 123,
				}).Return(nil)

				return authorDao, readerDao
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantId: 1,
		},
		{
			name: "新建帖子, 制作库成功, 线上库失败, 同步失败",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Create(gomock.Any(), dao.Article{
					Title:    "发表的标题",
					Content:  "发表的内容",
					AuthorId: 123,
				}).Return(int64(1), nil)

				readerDao.EXPECT().Upsert(gomock.Any(), dao.Article{
					Id:       1,
					Title:    "发表的标题",
					Content:  "发表的内容",
					AuthorId: 123,
				}).Return(errors.New("线上库失败"))

				return authorDao, readerDao
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr: errors.New("线上库失败"),
			wantId:  1,
		},
		{
			name: "新建帖子, 制作库失败, 同步失败",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Create(gomock.Any(), dao.Article{
					Title:    "发表的标题",
					Content:  "发表的内容",
					AuthorId: 123,
				}).Return(int64(1), errors.New("制作库失败"))

				return authorDao, readerDao
			},

			art: domain.Article{
				Title:   "发表的标题",
				Content: "发表的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr: errors.New("制作库失败"),
			wantId:  0,
		},
		{
			name: "修改帖子, 两库成功, 同步成功",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Update(gomock.Any(), dao.Article{
					Id:       2,
					Title:    "修改的标题",
					Content:  "修改的内容",
					AuthorId: 123,
				}).Return(nil)

				readerDao.EXPECT().Upsert(gomock.Any(), dao.Article{
					Id:       2,
					Title:    "修改的标题",
					Content:  "修改的内容",
					AuthorId: 123,
				}).Return(nil)

				return authorDao, readerDao
			},

			art: domain.Article{
				Id:      2,
				Title:   "修改的标题",
				Content: "修改的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantId: 2,
		},
		{
			name: "修改帖子, 制作库成功, 线上库失败, 同步失败",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Update(gomock.Any(), dao.Article{
					Id:       3,
					Title:    "修改的标题",
					Content:  "修改的内容",
					AuthorId: 123,
				}).Return(nil)

				readerDao.EXPECT().Upsert(gomock.Any(), dao.Article{
					Id:       3,
					Title:    "修改的标题",
					Content:  "修改的内容",
					AuthorId: 123,
				}).Return(errors.New("线上库失败"))

				return authorDao, readerDao
			},

			art: domain.Article{
				Id:      3,
				Title:   "修改的标题",
				Content: "修改的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr: errors.New("线上库失败"),
			wantId:  3,
		},
		{
			name: "修改帖子, 制作库失败, 同步失败",
			mock: func(ctrl *gomock.Controller) (dao.ArticleAuthorDao, dao.ArticleReaderDao) {
				authorDao := daomocks.NewMockArticleAuthorDao(ctrl)
				readerDao := daomocks.NewMockArticleReaderDao(ctrl)

				authorDao.EXPECT().Update(gomock.Any(), dao.Article{
					Id:       4,
					Title:    "修改的标题",
					Content:  "修改的内容",
					AuthorId: 123,
				}).Return(errors.New("制作库库失败"))

				return authorDao, readerDao
			},

			art: domain.Article{
				Id:      4,
				Title:   "修改的标题",
				Content: "修改的内容",
				Author: domain.Author{
					Id: 123,
				},
			},

			wantErr: errors.New("制作库库失败"),
			wantId:  4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authorDao, readerDao := tc.mock(ctrl)
			repo := NewCacheArticleRepositoryV2(authorDao, readerDao, logger.NewNopLogger())

			id, err := repo.Sync(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)

		})
	}
}
