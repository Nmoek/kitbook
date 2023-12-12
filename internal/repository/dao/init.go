package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},             //用户表
		&Article{},          //帖子表-制作库
		&PublishedArticle{}, //帖子表-线上库
		&Interactive{},      //互动表

	)
}

func InitCollection(mdb *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	colProduce := mdb.Collection("articles")
	_, err := colProduce.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true), // 唯一索引
		},
		{
			Keys: bson.D{{"author_id", 1}},
		},
	})

	colLive := mdb.Collection("published_articles")
	_, err = colLive.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true), // 唯一索引
		},
		{
			Keys: bson.D{{"author_id", 1}},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
