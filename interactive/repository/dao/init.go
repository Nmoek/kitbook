package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Interactive{},     //互动表, 阅读数+点赞数+收藏数
		&UserLikeInfo{},    //用户点赞信息表
		&UserCollectInfo{}, //用户收藏信息表
	)
}
