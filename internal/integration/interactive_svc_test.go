// Package integration
// @Description: 互动模块-集成测试
package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	startup2 "kitbook/internal/integration/startup"
	"kitbook/internal/repository/dao"
	"kitbook/internal/service"
	"testing"
	"time"
)

const (
	fieldReadCnt    = "read_cnt"
	fieldLikeCnt    = "like_cnt"
	fieldCollectCnt = "collect_cnt"
)

type InteractiveSvcSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
	svc service.InteractiveService
}

func (i *InteractiveSvcSuite) SetupSuite() {
	i.db = startup2.InitDB()
	i.rdb = startup2.InitRedis()
	i.svc = startup2.NewInteractiveService()

}

func (i *InteractiveSvcSuite) TearDownTest() {
	t := i.T()
	var err error

	//// 清空制作库
	//err = i.db.Exec("truncate table `articles`").Error
	//assert.NoError(t, err)
	//
	//// 清空线上库
	//err = i.db.Exec("truncate table `published_articles`").Error
	//assert.NoError(t, err)

	// 清空互动库
	err = i.db.Exec("truncate table `interactives`").Error
	assert.NoError(t, err)

	// 清空点赞信息库
	err = i.db.Exec("truncate table `user_like_infos`").Error
	assert.NoError(t, err)

	// 清空手收藏信息库
	err = i.db.Exec("truncate table `user_collect_infos`").Error
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 清空Redis
	err = i.rdb.FlushDB(ctx).Err()
	assert.NoError(t, err)

}

// @func: Test_IncreaseReadCnt
// @date: 2023-12-12 02:09:55
// @brief: 阅读数+1
// @author: Kewin Li
// @receiver a
func (i *InteractiveSvcSuite) TestIncreaseReadCnt() {
	t := i.T()

	testCases := []struct {
		name string

		before func(t *testing.T)

		after func(t *testing.T)

		biz   string
		bizId int64

		wantErr error
	}{
		// 新建记录,数据库保存成功,无缓存
		{
			name:   "New Record, Database Exists, Cache None",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 1. 查数据库
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 1, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 0)
				assert.True(t, intr.Ctime > 0)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:      1,
					BizId:   1,
					Biz:     "test",
					ReadCnt: 1,
				}, intr)

				// 2. 查缓存
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				key := createKey("test", 1)
				res, err := i.rdb.Exists(ctx, key).Result()
				assert.Equal(t, int64(0), res)

			},

			biz:   "test",
			bizId: 1,
		},
		// 修改记录, 数据库成功, 缓存成功
		{
			name: "Modify Record, Database Success, Cache Success",
			before: func(t *testing.T) {

				// 1. 数据库放数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
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
				err = i.rdb.HSet(ctx, createKey("test", 2), fieldReadCnt, 1).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 1. 查数据库
				var interactive dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 2, "test").First(&interactive).Error
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
				readCnt, err := i.rdb.HGet(ctx, key, fieldReadCnt).Int64()
				assert.NoError(t, err)
				assert.Equal(t, int64(2), readCnt)

				err = i.rdb.Del(ctx, key, fieldReadCnt).Err()
				assert.NoError(t, err)
			},

			biz:   "test",
			bizId: 2,
		},
		// 修改记录, 数据库成功, 缓存无
		{
			name: "Modify Record, Database Success, No Cache",
			before: func(t *testing.T) {

				// 1. 数据库放数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
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
				err := i.db.Where("biz_id = ? AND biz = ?", 3, "test").First(&interactive).Error
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
				res, err := i.rdb.Exists(ctx, key).Result()
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

			err := i.svc.IncreaseReadCnt(context.Background(), tc.biz, tc.bizId)

			assert.Equal(t, tc.wantErr, err)

		})

	}

}

// @func: Test_IncreaseReadCnt
// @date: 2024-01-06 21:04:40
// @brief: 点赞数+1
// @author: Kewin Li
// @receiver a
// @return func
func (i *InteractiveSvcSuite) TestLike() {
	t := i.T()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		bizId  int64
		biz    string
		userId int64
	}{
		// 新建记录, 数据库有, 缓存无
		{
			name:   "New Record, Database Exists, Cache None",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库 互动数据表
				err := i.db.Where("biz_id = ? AND biz = ? ", 1, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 0)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:      1,
					BizId:   1,
					Biz:     "test",
					LikeCnt: 1,
				}, intr)

				// 2. 查询数据库 点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 1, 1, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Utime > 0)
				assert.True(t, likeInfo.Ctime > 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Ctime = 0
				likeInfo.Utime = 0
				likeInfo.Id = 0
				assert.Equal(t, dao.UserLikeInfo{
					UserId: 1,
					BizId:  1,
					Biz:    "test",
					Status: 1,
				}, likeInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 1), fieldLikeCnt).Int()
				assert.Equal(t, redis.Nil, err)
				assert.Equal(t, 0, val)
			},

			bizId:  1,
			biz:    "test",
			userId: 1,
		},
		// 修改记录, 数据库成功, 缓存成功
		{
			name: "Modify Record, Database Success, Cache Success",
			before: func(t *testing.T) {
				// 1. 数据库预先放入数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					LikeCnt:    1,
					ReadCnt:    456,
					CollectCnt: 789,
					Utime:      666,
					Ctime:      666,
				}).Error
				assert.NoError(t, err)

				// 2. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 2), fieldLikeCnt, 1).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库
				err := i.db.Where("biz_id = ? AND biz = ? ", 2, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 666)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					LikeCnt:    2,
					ReadCnt:    456,
					CollectCnt: 789,
				}, intr)

				// 2. 查询数据库 点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 2, 2, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Utime > 0)
				assert.True(t, likeInfo.Ctime > 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Ctime = 0
				likeInfo.Utime = 0
				likeInfo.Id = 0
				assert.Equal(t, dao.UserLikeInfo{
					UserId: 2,
					BizId:  2,
					Biz:    "test",
					Status: 1,
				}, likeInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 2), fieldLikeCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, val)
			},

			bizId:  2,
			biz:    "test",
			userId: 2,
		},
		// 修改记录, 数据库成功, 缓存无
		{
			name: "Modify Record, Database Success, No Cache",
			before: func(t *testing.T) {
				// 1. 数据库预先放入数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					LikeCnt:    3,
					ReadCnt:    456,
					CollectCnt: 789,
					Utime:      777,
					Ctime:      777,
				}).Error
				assert.NoError(t, err)

				// 2. 缓存不放入
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库
				err := i.db.Where("biz_id = ? AND biz = ? ", 3, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 777)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					LikeCnt:    4, // 3 + 1
					ReadCnt:    456,
					CollectCnt: 789,
				}, intr)

				// 2. 查询数据库 点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 3, 3, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Utime > 0)
				assert.True(t, likeInfo.Ctime > 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Ctime = 0
				likeInfo.Utime = 0
				likeInfo.Id = 0
				assert.Equal(t, dao.UserLikeInfo{
					UserId: 3,
					BizId:  3,
					Biz:    "test",
					Status: 1,
				}, likeInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err = i.rdb.HGet(ctx, createKey("test", 3), fieldLikeCnt).Int()
				assert.Equal(t, redis.Nil, err)
			},

			bizId:  3,
			biz:    "test",
			userId: 3,
		},
		// TODO: 同一个用户重复点赞同一个内容, 不知道怎么实现代码

		//{
		//	name: "Repeatedly clicking in favor of a piece of content by the same user",
		//	before: func(t *testing.T) {
		//		// 1. 数据库放入数据-互动数据表
		//		err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
		//			Id:         4,
		//			BizId:      4,
		//			Biz:        "test",
		//			LikeCnt:    4,
		//			ReadCnt:    456,
		//			CollectCnt: 789,
		//			Utime:      777,
		//			Ctime:      777,
		//		}).Error
		//		assert.NoError(t, err)
		//
		//		// 2. 数据库放入数据-点赞信息表
		//		err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
		//			Id:    4,
		//			BizId: 4,
		//			Biz:   "test",
		//			Utime: 888,
		//			Ctime: 888,
		//		}).Error
		//		assert.NoError(t, err)
		//		// 3. 缓存放入
		//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		//		defer cancel()
		//		err = i.rdb.HSet(ctx, createKey("test", 4), fieldLikeCnt, 4).Err()
		//		assert.NoError(t, err)
		//	},
		//	after: func(t *testing.T) {
		//		var intr dao.Interactive
		//		// 1. 查询数据库
		//		err := i.db.Where("biz_id = ? AND biz = ? ", 4, "test").First(&intr).Error
		//		assert.NoError(t, err)
		//		assert.True(t, intr.Ctime > 0)
		//		assert.True(t, intr.Utime == 777)
		//		intr.Ctime = 0
		//		intr.Utime = 0
		//		assert.Equal(t, dao.Interactive{
		//			Id:         4,
		//			BizId:      4,
		//			Biz:        "test",
		//			LikeCnt:    4,
		//			ReadCnt:    456,
		//			CollectCnt: 789,
		//		}, intr)
		//
		//		// 2. 查询数据库 点赞信息表
		//		var likeInfo dao.UserLikeInfo
		//		err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 3, 3, "test").First(&likeInfo).Error
		//		assert.NoError(t, err)
		//		assert.True(t, likeInfo.Utime == 888)
		//		assert.True(t, likeInfo.Ctime > 0)
		//		assert.True(t, likeInfo.Status == 1)
		//		likeInfo.Ctime = 0
		//		likeInfo.Utime = 0
		//		assert.Equal(t, dao.UserLikeInfo{
		//			Id:     1,
		//			UserId: 3,
		//			BizId:  3,
		//			Biz:    "test",
		//			Status: 1,
		//		}, likeInfo)
		//
		//		// 3. 查询缓存
		//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		//		defer cancel()
		//		val, err := i.rdb.HGet(ctx, createKey("test", 4), fieldLikeCnt).Int()
		//		assert.NoError(t, err)
		//		assert.Equal(t, 4, val)
		//	},
		//
		//	bizId:  4,
		//	biz:    "test",
		//	userId: 4,
		//},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			err := i.svc.Like(context.Background(), tc.biz, tc.bizId, tc.userId)
			assert.NoError(t, err)

		})

	}

}

// @func: TestCancelLike
// @date: 2024-01-07 00:22:42
// @brief: 点赞数-1
// @author: Kewin Li
// @receiver i
func (i *InteractiveSvcSuite) TestCancelLike() {
	t := i.T()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		userId int64
		bizId  int64
		biz    string

		wantErr error
	}{
		// 修改记录，数据库成功，缓存成功
		{
			name: "Modify Record, Database Success, Cache Success",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         1,
					BizId:      1,
					Biz:        "test",
					LikeCnt:    1,
					ReadCnt:    666,
					CollectCnt: 666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 点赞信息表
				err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     1,
					UserId: 1,
					BizId:  1,
					Biz:    "test",
					Status: 1, // 点赞有效
					Utime:  999,
					Ctime:  999,
				}).Error
				assert.NoError(t, err)

				// 3. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 1), fieldLikeCnt, 1).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 1, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime > 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					BizId:      1,
					Biz:        "test",
					LikeCnt:    0,
					ReadCnt:    666,
					CollectCnt: 666,
				}, intr)

				// 2. 数据库查询-点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 1, 1, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Ctime == 999)
				assert.True(t, likeInfo.Utime > 999)
				assert.True(t, likeInfo.Status == 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Utime = 0
				likeInfo.Ctime = 0
				likeInfo.Id = 0

				assert.Equal(t, dao.UserLikeInfo{
					UserId: 1,
					BizId:  1,
					Biz:    "test",
					Status: 0,
				}, likeInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 1), fieldLikeCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 0, val)
			},

			userId: 1,
			bizId:  1,
			biz:    "test",
		},
		// 修改记录，数据库成功，缓存失败
		{
			name: "Modify Record, Database Success, Cache Failure",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					LikeCnt:    2,
					ReadCnt:    666,
					CollectCnt: 666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 点赞信息表
				err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     2,
					UserId: 2,
					BizId:  2,
					Biz:    "test",
					Status: 1, // 点赞有效
					Utime:  999,
					Ctime:  999,
				}).Error
				assert.NoError(t, err)

				// 3. 缓存不放入

			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 2, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime > 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					LikeCnt:    1,
					ReadCnt:    666,
					CollectCnt: 666,
				}, intr)

				// 2. 数据库查询-点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 2, 2, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Ctime == 999)
				assert.True(t, likeInfo.Utime > 999)
				assert.True(t, likeInfo.Status == 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Utime = 0
				likeInfo.Ctime = 0
				likeInfo.Id = 0

				assert.Equal(t, dao.UserLikeInfo{
					UserId: 2,
					BizId:  2,
					Biz:    "test",
					Status: 0,
				}, likeInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err = i.rdb.HGet(ctx, createKey("test", 2), fieldLikeCnt).Int()
				assert.Equal(t, redis.Nil, err)
			},

			userId: 2,
			bizId:  2,
			biz:    "test",
		},
		// 边界下界测试:点赞数为0, 数据库失败, 缓存不更新
		{
			name: "Lower Bound Test: Likes 0, Database Failure, Cache Not Updated",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					LikeCnt:    0,
					ReadCnt:    666,
					CollectCnt: 666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 点赞信息表
				i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     3,
					UserId: 3,
					BizId:  3,
					Biz:    "test",
					Status: 0,
					Utime:  999,
					Ctime:  999,
				})

				// 3. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 3), fieldLikeCnt, 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 3, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime == 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					LikeCnt:    0,
					ReadCnt:    666,
					CollectCnt: 666,
				}, intr)

				// 2. 数据库查询-点赞信息表
				var likeInfo dao.UserLikeInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 3, 3, "test").First(&likeInfo).Error
				assert.NoError(t, err)
				assert.True(t, likeInfo.Ctime == 999)
				assert.True(t, likeInfo.Utime == 999)
				assert.True(t, likeInfo.Status == 0)
				assert.True(t, likeInfo.Id > 0)
				likeInfo.Utime = 0
				likeInfo.Ctime = 0
				likeInfo.Id = 0

				assert.Equal(t, dao.UserLikeInfo{
					UserId: 3,
					BizId:  3,
					Biz:    "test",
					Status: 0,
				}, likeInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 3), fieldLikeCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 0, val)
			},

			userId:  3,
			bizId:   3,
			biz:     "test",
			wantErr: errors.New("用户非法操作"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			err := i.svc.CancelLike(context.Background(), tc.biz, tc.bizId, tc.userId)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

// @func: TestCollect
// @date: 2024-01-07 17:47:38
// @brief: 收藏数+1
// @author: Kewin Li
// @receiver i
func (i *InteractiveSvcSuite) TestCollect() {
	t := i.T()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		userId    int64
		bizId     int64
		biz       string
		collectId int64

		wantErr error
	}{
		// 新建记录, 数据库有, 缓存无
		{
			name:   "New Record, Database Exists, Cache None",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库 互动数据表
				err := i.db.Where("biz_id = ? AND biz = ? ", 1, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 0)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:         1,
					BizId:      1,
					Biz:        "test",
					CollectCnt: 1,
				}, intr)

				// 2. 查询数据库 收藏信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 1, 1, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Utime > 0)
				assert.True(t, collectInfo.Ctime > 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Ctime = 0
				collectInfo.Utime = 0
				collectInfo.Id = 0
				assert.Equal(t, dao.UserCollectInfo{
					UserId:    1,
					BizId:     1,
					Biz:       "test",
					CollectId: 1, //指定收藏到哪个收藏夹
					Status:    1,
				}, collectInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err = i.rdb.HGet(ctx, createKey("test", 1), fieldCollectCnt).Int()
				assert.Equal(t, redis.Nil, err)
			},

			bizId:     1,
			biz:       "test",
			userId:    1,
			collectId: 1,
		},
		// 修改记录, 数据库成功, 缓存成功
		{
			name: "Modify Record, Database Success, Cache Success",
			before: func(t *testing.T) {
				// 1. 数据库预先放入数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					CollectCnt: 1,
					ReadCnt:    456,
					LikeCnt:    789,
					Utime:      666,
					Ctime:      666,
				}).Error
				assert.NoError(t, err)

				// 2. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 2), fieldCollectCnt, 1).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库
				err := i.db.Where("biz_id = ? AND biz = ? ", 2, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 666)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					CollectCnt: 2,
					ReadCnt:    456,
					LikeCnt:    789,
				}, intr)

				// 2. 查询数据库 收藏信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 2, 2, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Utime > 0)
				assert.True(t, collectInfo.Ctime > 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Ctime = 0
				collectInfo.Utime = 0
				collectInfo.Id = 0
				assert.Equal(t, dao.UserCollectInfo{
					UserId:    2,
					BizId:     2,
					Biz:       "test",
					CollectId: 2,
					Status:    1,
				}, collectInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 2), fieldCollectCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, val)
			},

			bizId:     2,
			biz:       "test",
			userId:    2,
			collectId: 2,
		},
		// 修改记录, 数据库成功, 缓存无
		{
			name: "Modify Record, Database Success, No Cache",
			before: func(t *testing.T) {
				// 1. 数据库预先放入数据
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					CollectCnt: 3,
					ReadCnt:    456,
					LikeCnt:    789,
					Utime:      777,
					Ctime:      777,
				}).Error
				assert.NoError(t, err)

				// 2. 缓存不放入
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				// 1. 查询数据库
				err := i.db.Where("biz_id = ? AND biz = ? ", 3, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				assert.True(t, intr.Utime > 777)
				intr.Ctime = 0
				intr.Utime = 0

				assert.Equal(t, dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					CollectCnt: 4, // 3 + 1
					ReadCnt:    456,
					LikeCnt:    789,
				}, intr)

				// 2. 查询数据库 收藏信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 3, 3, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Utime > 0)
				assert.True(t, collectInfo.Ctime > 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Ctime = 0
				collectInfo.Utime = 0
				collectInfo.Id = 0
				assert.Equal(t, dao.UserCollectInfo{
					UserId:    3,
					BizId:     3,
					Biz:       "test",
					CollectId: 3,
					Status:    1,
				}, collectInfo)

				// 3. 查询缓存
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err = i.rdb.HGet(ctx, createKey("test", 3), fieldCollectCnt).Int()
				assert.Equal(t, redis.Nil, err)
			},

			bizId:     3,
			biz:       "test",
			userId:    3,
			collectId: 3,
		},
		// TODO: 同一个用户重复收藏同一个内容, 不知道怎么实现代码
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			err := i.svc.Collect(context.Background(), tc.biz, tc.bizId, tc.collectId, tc.userId)
			assert.Equal(t, tc.wantErr, err)
		})
	}

}

// @func: TestCancelCollect
// @date: 2024-01-07 17:50:40
// @brief: 收藏数-1
// @author: Kewin Li
// @receiver i
func (i *InteractiveSvcSuite) TestCancelCollect() {
	t := i.T()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		userId    int64
		bizId     int64
		biz       string
		collectId int64

		wantErr error
	}{
		// 修改记录，数据库成功，缓存成功
		{
			name: "Modify Record, Database Success, Cache Success",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         1,
					BizId:      1,
					Biz:        "test",
					CollectCnt: 1,
					ReadCnt:    666,
					LikeCnt:    666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 点赞信息表
				err = i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:        1,
					UserId:    1,
					BizId:     1,
					Biz:       "test",
					Status:    1, // 收藏有效
					CollectId: 1,
					Utime:     999,
					Ctime:     999,
				}).Error
				assert.NoError(t, err)

				// 3. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 1), fieldCollectCnt, 1).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 1, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime > 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         1,
					BizId:      1,
					Biz:        "test",
					CollectCnt: 0,
					ReadCnt:    666,
					LikeCnt:    666,
				}, intr)

				// 2. 数据库查询-点赞信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 1, 1, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Ctime == 999)
				assert.True(t, collectInfo.Utime > 999)
				assert.True(t, collectInfo.Status == 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Utime = 0
				collectInfo.Ctime = 0
				collectInfo.Id = 0

				assert.Equal(t, dao.UserCollectInfo{
					UserId:    1,
					BizId:     1,
					Biz:       "test",
					CollectId: 1,
					Status:    0,
				}, collectInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 1), fieldCollectCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 0, val)
			},

			userId:    1,
			bizId:     1,
			collectId: 1,
			biz:       "test",
		},
		// 修改记录，数据库成功，缓存失败
		{
			name: "Modify Record, Database Success, Cache Failure",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					CollectCnt: 2,
					ReadCnt:    666,
					LikeCnt:    666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 收藏信息表
				err = i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:        2,
					UserId:    2,
					BizId:     2,
					Biz:       "test",
					Status:    1, // 点赞有效
					CollectId: 2,
					Utime:     999,
					Ctime:     999,
				}).Error
				assert.NoError(t, err)

				// 3. 缓存不放入

			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 2, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime > 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         2,
					BizId:      2,
					Biz:        "test",
					CollectCnt: 1,
					ReadCnt:    666,
					LikeCnt:    666,
				}, intr)

				// 2. 数据库查询-收藏信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 2, 2, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Ctime == 999)
				assert.True(t, collectInfo.Utime > 999)
				assert.True(t, collectInfo.Status == 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Utime = 0
				collectInfo.Ctime = 0
				collectInfo.Id = 0

				assert.Equal(t, dao.UserCollectInfo{
					UserId:    2,
					BizId:     2,
					Biz:       "test",
					CollectId: 2,
					Status:    0,
				}, collectInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err = i.rdb.HGet(ctx, createKey("test", 2), fieldCollectCnt).Int()
				assert.Equal(t, redis.Nil, err)
			},

			userId:    2,
			bizId:     2,
			collectId: 2,
			biz:       "test",
		},
		// 边界下界测试:收藏数为0, 数据库失败, 缓存不更新
		{
			name: "Lower Bound Test: Likes 0, Database Failure, Cache Not Updated",
			before: func(t *testing.T) {
				// 1. 数据库放入数据- 互动数据表
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					CollectCnt: 0,
					ReadCnt:    666,
					LikeCnt:    666,
					Utime:      888,
					Ctime:      888,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入数据- 收藏信息表
				i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:        3,
					UserId:    3,
					BizId:     3,
					Biz:       "test",
					Status:    0,
					CollectId: 3,
					Utime:     999,
					Ctime:     999,
				})

				// 3. 缓存放入
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = i.rdb.HSet(ctx, createKey("test", 3), fieldLikeCnt, 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {

				// 1. 数据库查询-互动数据表
				var intr dao.Interactive
				err := i.db.Where("biz_id = ? AND biz = ?", 3, "test").First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime == 888)
				assert.True(t, intr.Utime == 888)
				intr.Utime = 0
				intr.Ctime = 0
				assert.Equal(t, dao.Interactive{
					Id:         3,
					BizId:      3,
					Biz:        "test",
					CollectCnt: 0,
					ReadCnt:    666,
					LikeCnt:    666,
				}, intr)

				// 2. 数据库查询-点赞信息表
				var collectInfo dao.UserCollectInfo
				err = i.db.Where("user_id = ? AND biz_id = ? AND biz = ?", 3, 3, "test").First(&collectInfo).Error
				assert.NoError(t, err)
				assert.True(t, collectInfo.Ctime == 999)
				assert.True(t, collectInfo.Utime == 999)
				assert.True(t, collectInfo.Status == 0)
				assert.True(t, collectInfo.Id > 0)
				collectInfo.Utime = 0
				collectInfo.Ctime = 0
				collectInfo.Id = 0

				assert.Equal(t, dao.UserCollectInfo{
					UserId:    3,
					BizId:     3,
					Biz:       "test",
					CollectId: 3,
					Status:    0,
				}, collectInfo)

				// 3. 缓存查询
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := i.rdb.HGet(ctx, createKey("test", 3), fieldLikeCnt).Int()
				assert.NoError(t, err)
				assert.Equal(t, 0, val)
			},

			userId:    3,
			bizId:     3,
			collectId: 3,
			biz:       "test",
			wantErr:   errors.New("用户非法操作"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			err := i.svc.CancelCollect(context.Background(), tc.biz, tc.bizId, tc.collectId, tc.userId)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// @func: TestGet
// @date: 2024-01-07 19:01:08
// @brief: 聚合查询-阅读数/点赞数/收藏数
// @author: Kewin Li
// @receiver i
func (i *InteractiveSvcSuite) TestGet() {
	t := i.T()
	now := time.Now().UnixMilli()

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		biz    string
		bizId  int64
		userId int64

		wantRes domain.Interactive
		wantErr error
	}{
		// 查询缓存,点赞信息表有记录,收藏信息表有记录
		{
			name: "Query cache, likes are recorded, collects are recorded",
			before: func(t *testing.T) {
				// 1. 缓存放入 阅读数/点赞数/收藏数
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := i.rdb.HSet(ctx, createKey("test", 1), map[string]any{
					"biz_id":        1,
					fieldReadCnt:    111,
					fieldLikeCnt:    222,
					fieldCollectCnt: 333,
				}).Err()
				assert.NoError(t, err)

				// 2. 数据库放入 点赞信息表
				err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     1,
					UserId: 1,
					BizId:  1,
					Biz:    "test",
					Status: 1,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)

				// 3. 数据库放入 收藏信息表
				err = i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:     1,
					UserId: 1,
					BizId:  1,
					Biz:    "test",
					Status: 1,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {},

			biz:    "test",
			bizId:  1,
			userId: 1,

			wantRes: domain.Interactive{
				BizId:      1,
				ReadCnt:    111,
				LikeCnt:    222,
				CollectCnt: 333,
				Liked:      true,
				Collected:  true,
			},
		},
		// 查询缓存,点赞信息表取消点赞,收藏信息表取消收藏
		{
			name: "Query Cache,Likes Unliked,Collections UnCollected",
			before: func(t *testing.T) {
				// 1. 缓存放入 阅读数/点赞数/收藏数
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := i.rdb.HSet(ctx, createKey("test", 2), map[string]any{
					"biz_id":        2,
					fieldReadCnt:    111,
					fieldLikeCnt:    222,
					fieldCollectCnt: 333,
				}).Err()
				assert.NoError(t, err)

				// 2. 数据库放入 点赞信息表
				err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     2,
					UserId: 2,
					BizId:  2,
					Biz:    "test",
					Status: 0,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)
				// 3. 数据库放入 收藏信息表
				err = i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:     2,
					UserId: 2,
					BizId:  2,
					Biz:    "test",
					Status: 0,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {},

			biz:    "test",
			bizId:  2,
			userId: 2,

			wantRes: domain.Interactive{
				BizId:      2,
				ReadCnt:    111,
				LikeCnt:    222,
				CollectCnt: 333,
				Liked:      false,
				Collected:  false,
			},
		},
		// 查询缓存,点赞信息表无记录,收藏信息表无记录
		{
			name: "Query Cache, No record for Like table, No record for Favorite table",
			before: func(t *testing.T) {
				// 1. 缓存放入 阅读数/点赞数/收藏数
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := i.rdb.HSet(ctx, createKey("test", 3), map[string]any{
					"biz_id":        3,
					fieldReadCnt:    111,
					fieldLikeCnt:    222,
					fieldCollectCnt: 333,
				}).Err()
				assert.NoError(t, err)

				// 2. 数据库不放入 点赞信息表

				// 3. 数据库不放入 收藏信息表

			},
			after: func(t *testing.T) {},

			biz:    "test",
			bizId:  3,
			userId: 3,

			wantRes: domain.Interactive{
				BizId:      3,
				ReadCnt:    111,
				LikeCnt:    222,
				CollectCnt: 333,
				Liked:      false,
				Collected:  false,
			},
		},
		// 查询库,缓存回写成功,点赞信息表有记录,收藏信息表有记录
		{
			name: "Query library, cache writeback successful, likes table has records, favorites table has records",
			before: func(t *testing.T) {
				// 1. 数据库放入 阅读数/点赞数/收藏数
				err := i.db.Model(&dao.Interactive{}).Create(dao.Interactive{
					Id:         4,
					BizId:      4,
					Biz:        "test",
					ReadCnt:    111,
					LikeCnt:    222,
					CollectCnt: 333,
					Utime:      now,
					Ctime:      now,
				}).Error
				assert.NoError(t, err)

				// 2. 数据库放入 点赞信息表
				err = i.db.Model(&dao.UserLikeInfo{}).Create(dao.UserLikeInfo{
					Id:     4,
					UserId: 4,
					BizId:  4,
					Biz:    "test",
					Status: 1,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)

				// 3. 数据库放入 收藏信息表
				err = i.db.Model(&dao.UserCollectInfo{}).Create(dao.UserCollectInfo{
					Id:     4,
					UserId: 4,
					BizId:  4,
					Biz:    "test",
					Status: 1,
					Utime:  now,
					Ctime:  now,
				}).Error
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {},

			biz:    "test",
			bizId:  4,
			userId: 4,

			wantRes: domain.Interactive{
				BizId:      4,
				ReadCnt:    111,
				LikeCnt:    222,
				CollectCnt: 333,
				Liked:      true,
				Collected:  true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			start := time.Now()
			res, err := i.svc.Get(context.Background(), tc.biz, tc.bizId, tc.userId)
			cost := time.Since(start).Milliseconds()
			t.Log("cost=", cost)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

		})
	}

}

// @func: TestGetByIds
// @date: 2024-01-07 22:49:57
// @brief: 批量查询-互动表数据
// @author: Kewin Li
// @receiver i
func (i *InteractiveSvcSuite) TestGetByIds() {
	t := i.T()
	now := time.Now().UnixMilli()
	intrs := make([]dao.Interactive, 10)

	testCases := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		biz    string
		bizIds []int64

		wantRes map[int64]domain.Interactive
		wantErr error
	}{
		// 数据库没有记录
		{
			name:   "No records in the database",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},

			biz:     "test",
			bizIds:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantRes: map[int64]domain.Interactive{},
		},
		// 批量查询数据库记录
		{
			name: "Batch query database records",
			before: func(t *testing.T) {
				for i := 0; i < len(intrs); i++ {
					intrs[i].Id = int64(i) + 1
					intrs[i].BizId = int64(i) + 1
					intrs[i].Biz = "test"
					intrs[i].ReadCnt = 111
					intrs[i].LikeCnt = 222
					intrs[i].CollectCnt = 333
					intrs[i].Ctime = now
					intrs[i].Utime = now
				}

				err := i.db.Model(&dao.Interactive{}).CreateInBatches(&intrs, 10).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {},

			biz:    "test",
			bizIds: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			wantRes: map[int64]domain.Interactive{
				1:  {BizId: 1, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				2:  {BizId: 2, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				3:  {BizId: 3, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				4:  {BizId: 4, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				5:  {BizId: 5, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				6:  {BizId: 6, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				7:  {BizId: 7, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				8:  {BizId: 8, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				9:  {BizId: 9, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
				10: {BizId: 10, ReadCnt: 111, LikeCnt: 222, CollectCnt: 333},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			res, err := i.svc.GetByIds(context.Background(), tc.biz, tc.bizIds)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRes, res)

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
