package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	domain2 "kitbook/interactive/domain"
	"kitbook/interactive/service"
	"kitbook/internal/domain"
	svcmocks "kitbook/internal/service/mocks"
	"testing"
	"time"
)

// @func: TestBatchRankingService_topN
// @date: 2023-12-27 00:31:04
// @brief: 测试热榜功能-核心算法实现
// @author: Kewin Li
// @param t
func TestBatchRankingService_topN(t *testing.T) {

	const batchSize = 2 // 一批获取2条记录
	utime := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService)

		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "成功获取",
			mock: func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				/*查询帖子内容*/
				// 获取第一批记录
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 2).Return([]domain.Article{
					{
						Id:    1,
						Utime: utime,
					},
					{
						Id:    2,
						Utime: utime,
					},
				}, nil)

				// 获取第二批记录
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 2, 2).Return([]domain.Article{
					{
						Id:    3,
						Utime: utime,
					},
					{
						Id:    4,
						Utime: utime,
					},
				}, nil)
				// 获取第三批记录
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 4, 2).Return([]domain.Article{}, nil)

				/*查询互动数据*/
				// 获取第一批数据
				intrSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2}).Return(map[int64]domain2.Interactive{
					1: {LikeCnt: 1},
					2: {LikeCnt: 2},
				}, nil)

				// 获取第二批数据
				intrSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{3, 4}).Return(map[int64]domain2.Interactive{
					3: {LikeCnt: 3},
					4: {LikeCnt: 4},
				}, nil)

				// 获取第三批数据
				//intrSvc.EXPECT().GetByIds(gomock.Any(), "article", nil).Return(map[int64]domain.Interactive{}, nil)

				return intrSvc, artSvc
			},

			wantErr: nil,
			wantArts: []domain.Article{
				{Id: 4, Utime: utime},
				{Id: 3, Utime: utime},
				{Id: 2, Utime: utime},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 查找点赞数、 查找文章
			intrSvc, artSvc := tc.mock(ctrl)

			// 没有必要使用真实生产环境下的数量和计算函数
			svc := &BatchRankingService{
				intrSvc:   intrSvc,
				artSvc:    artSvc,
				batchSize: batchSize,
				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
					return float64(likeCnt)
				},
				n: 3,
			}

			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArts, arts)

		})
	}
}
