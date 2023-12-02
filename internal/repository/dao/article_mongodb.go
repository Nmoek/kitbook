package dao

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBArticleDAO struct {
	node       *snowflake.Node
	produceCol *mongo.Collection
	liveCol    *mongo.Collection
}

func NewMongoDBArticleDAO(mdb *mongo.Database, node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		node:       node,
		produceCol: mdb.Collection("articles"),
		liveCol:    mdb.Collection("published_articles"),
	}
}

// @func: Insert
// @date: 2023-12-01 01:41:09
// @brief: mongodb-插入文档
// @author: Kewin Li
// @receiver m
// @param ctx
// @param art
// @return int64
// @return error
func (m *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()

	art.Id = m.node.Generate().Int64()
	art.Ctime = now
	art.Utime = now
	_, err := m.produceCol.InsertOne(ctx, &art)
	return art.Id, err

}

// @func: UpdateById
// @date: 2023-12-01 02:26:16
// @brief: mongodb-更新文档按ID
// @author: Kewin Li
// @receiver m
// @param ctx
// @param art
// @return error
func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, art Article) error {
	updateFilter := bson.M{
		"id":        art.Id,
		"author_id": art.AuthorId,
	}

	updateRes, err := m.produceCol.UpdateOne(ctx, updateFilter, bson.D{{"$set", bson.M{
		"title":     art.Title,
		"content":   art.Content,
		"author_id": art.AuthorId,
		"status":    art.Status,
		"utime":     time.Now().UnixMilli(),
	}}})
	if updateRes.ModifiedCount <= 0 {
		return ErrUserMismatch
	}
	return err
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {

	var err error

	if art.Id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		art.Id, err = m.Insert(ctx, art)
	}

	if err != nil {
		return art.Id, err
	}

	/*操作库*/
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.D{
		{"id", art.Id},
		{"author_id", art.AuthorId},
	}

	set := bson.D{
		{"$set", PublishedArticle(art)},
		{"$setOnInsert", bson.D{
			{"ctime", now},
		}},
	}

	_, err = m.liveCol.UpdateOne(ctx, filter, set,
		options.Update().SetUpsert(true))

	// 线上库遵循啊UPSERT语义
	return art.Id, err
}

func (m *MongoDBArticleDAO) SyncStatus(ctx *gin.Context, artId int64, authorId int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
