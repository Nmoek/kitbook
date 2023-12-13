// Package integration
// @Description: 互动模块-集成测试
package integration

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"kitbook/integration/startup"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"testing"
	"time"
)

const fieldReadCnt = "read_cnt"

type InteractiveSvcSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
	svc service.InteractiveService
}

func (a *InteractiveSvcSuite) SetupSuite() {
	a.db = startup.InitDB()
	a.rdb = startup.InitRedis()
}

func (a *InteractiveSvcSuite) TearDownTest() {
	a.db.Exec("truncate table `articles`")
	a.db.Exec("truncate table `published_articles`")
	a.db.Exec("truncate table `interactives`")
}

// @func: Test_IncreaseReadCnt
// @date: 2023-12-12 02:09:55
// @brief: 阅读数+1
// @author: Kewin Li
// @receiver a
func (a *InteractiveSvcSuite) Test_IncreaseReadCnt() {
	t := a.T()
	a.svc = startup.NewInteractiveService()

	testCases := []struct {
		name string

		before func(t *testing.T)

		after func(t *testing.T)

		biz   string
		bizId int64

		wantErr error
	}{
		{
			name:   "新建记录, 数据库保存成功, 无缓存",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 1. 查数据库
				var interactive dao.Interactive
				err := a.db.Where("biz_id = ?", 1).First(&interactive).Error
				assert.NoError(t, err)
				assert.True(t, interactive.Utime > 0)
				assert.True(t, interactive.Ctime > 0)
				interactive.Utime = 0
				interactive.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:      1,
					BizId:   1,
					Biz:     "test",
					ReadCnt: 1,
				}, interactive)

				// 2. 查缓存
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				key := createKey("test", 1)
				res, err := a.rdb.Exists(ctx, key).Result()
				assert.Equal(t, int64(0), res)

			},

			biz:   "test",
			bizId: 1,
		},
		{
			name: "数据库修改成功, 缓存更新成功",
			before: func(t *testing.T) {

				// 1. 数据库放数据
				err := a.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					ReadCnt:    1,
					LikeCnt:    1,
					CollectCnt: 2,
					Utime:      666,
					Ctime:      666,
				}).Error
				assert.NoError(t, err)

				// 2. 缓存放数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = a.rdb.HSet(ctx, createKey("test", 2), fieldReadCnt, 1).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 1. 查数据库
				var interactive dao.Interactive
				err := a.db.Where("biz_id = ?", 2).First(&interactive).Error
				assert.NoError(t, err)
				assert.True(t, interactive.Utime > 0)
				assert.True(t, interactive.Ctime == 666)
				interactive.Utime = 0
				interactive.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					ReadCnt:    2,
					LikeCnt:    1,
					CollectCnt: 2,
				}, interactive)

				// 2. 查缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				key := createKey("test", 2)
				readCnt, err := a.rdb.HGet(ctx, key, fieldReadCnt).Int64()
				assert.NoError(t, err)
				assert.Equal(t, int64(2), readCnt)

				err = a.rdb.Del(ctx, key, fieldReadCnt).Err()
				assert.NoError(t, err)
			},

			biz:   "test",
			bizId: 2,
		},
		{
			name: "数据库修改成功, 无缓存",
			before: func(t *testing.T) {

				// 1. 数据库放数据
				err := a.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					ReadCnt:    3,
					LikeCnt:    1,
					CollectCnt: 2,
					Utime:      777,
					Ctime:      777,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				// 1. 查数据库
				var interactive dao.Interactive
				err := a.db.Model(&dao.Interactive{}).Where("biz_id = ?", 3).First(&interactive).Error
				assert.NoError(t, err)
				assert.True(t, interactive.Utime > 0)
				assert.True(t, interactive.Ctime == 777)
				interactive.Utime = 0
				interactive.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					ReadCnt:    4, // 阅读数
					LikeCnt:    1,
					CollectCnt: 2,
				}, interactive)

				// 2. 查缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				key := createKey("test", 3)
				res, err := a.rdb.Exists(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), res)

			},

			biz:   "test",
			bizId: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			tc.before(t)
			defer tc.after(t)

			err := a.svc.IncreaseReadCnt(context.Background(), tc.biz, tc.bizId)

			assert.Equal(t, tc.wantErr, err)

		})

	}

}

// @func: TestInteractiveService
// @date: 2023-12-12 12:50:28
// @brief: 互动模块-测试套件入口
// @author: Kewin Li
// @param t
func TestInteractiveSvc(t *testing.T) {
	suite.Run(t, &InteractiveSvcSuite{})
}

// @func: createKey
// @date: 2023-12-12 23:53:51
// @brief: 制作interactive缓存的key
// @author: Kewin Li
// @receiver i
// @param biz
// @param bizId
// @return string
func createKey(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
