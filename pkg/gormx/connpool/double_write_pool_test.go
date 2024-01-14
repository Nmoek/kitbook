package connpool

import (
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitbook/pkg/logger"
	"testing"
	"time"
)

type DoubleWritePoolSuite struct {
	suite.Suite
	db  *gorm.DB
	src *gorm.DB
	dst *gorm.DB
}

func (d *DoubleWritePoolSuite) SetupSuite() {
	t := d.T()

	// 初始化源库
	src, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/kitbook"))
	assert.NoError(t, err)
	err = src.AutoMigrate(&Interactive{})
	d.src = src
	// 初始化目标库
	dst, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/kitbook_intr"))
	assert.NoError(t, err)
	err = dst.AutoMigrate(&Interactive{})
	assert.NoError(t, err)
	d.dst = dst

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: &DoubleWritePool{
			src:     src.ConnPool,
			dst:     dst.ConnPool,
			pattern: atomicx.NewValueOf(PatternSrcFirst),
			l:       logger.NewNopLogger(),
		},
	}))
	assert.NoError(t, err)
	d.db = db

}

func (d *DoubleWritePoolSuite) TearDownTest() {
	// 两个库中同一张表都会被清空
	// 有问题
	d.db.Exec("TRUNCATE TABLE `interactives`")
}

// @func: TestDoubleWritePool_ExecContext
// @date: 2024-01-14 19:10:55
// @brief: 双写-写入测试
// @author: Kewin Li
// @receiver d
func (d *DoubleWritePoolSuite) TestDoubleWritePool_Write() {
	t := d.T()
	now := time.Now().UnixMilli()
	testIntr := Interactive{
		BizId:      1,
		Biz:        "test",
		ReadCnt:    111,
		LikeCnt:    222,
		CollectCnt: 333,
		Ctime:      now,
		Utime:      now,
	}

	// TODO: 分情况进行双写

	err := d.db.Model(&Interactive{}).Create(&testIntr).Error
	assert.NoError(t, err)

	var srcIntr Interactive
	err = d.src.Where("biz_id = ?", 1).First(&srcIntr).Error
	assert.NoError(t, err)
	assert.Equal(t, testIntr, srcIntr)

	var dstIntr Interactive
	err = d.src.Where("biz_id = ?", 1).First(&dstIntr).Error
	assert.NoError(t, err)
	assert.Equal(t, testIntr, dstIntr)

}

func TestDoubleWritePool(t *testing.T) {
	suite.Run(t, &DoubleWritePoolSuite{})
}

// Interactive
// @Description: 阅读数、点赞数、收藏数三合一
type Interactive struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`

	// 建立联合唯一索引<bizId, biz
	BizId int64  `gorm:"uniqueIndex:intr_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:intr_biz_type_id"`

	// 阅读数
	ReadCnt int64
	// 点赞数
	LikeCnt int64
	// 收藏数
	CollectCnt int64
	Utime      int64
	Ctime      int64
}
