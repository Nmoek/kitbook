//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"kitbook/comment/grpc"
	"kitbook/comment/repository"
	"kitbook/comment/repository/dao"
	"kitbook/comment/service"
)

var thirdPartySet = wire.NewSet(
	InitDB,
	InitRedis,
	InitLogger,
	//InitSaramaClient,
	//InitSyncProducer,
	//InitConsumers,
)

var commentSvcSet = wire.NewSet(
	dao.NewGormCommentDao,
	repository.NewArticleCommentRepository,
	service.NewArticlePaymentService,
)

func NewCommentServiceServer() *grpc.ArticleCommentServiceServer {
	wire.Build(
		thirdPartySet,
		commentSvcSet,
		grpc.NewArticleCommentServiceServer,
	)
	return new(grpc.ArticleCommentServiceServer)
}
