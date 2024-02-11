//go:build wireinject

package main

import (
	"github.com/google/wire"
	"kitbook/comment/grpc"
	"kitbook/comment/ioc"
	"kitbook/comment/repository"
	"kitbook/comment/repository/dao"
	"kitbook/comment/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,

	//InitConsumers,
)

var commentSvcSet = wire.NewSet(
	dao.NewGormCommentDao,
	repository.NewArticleCommentRepository,
	service.NewArticlePaymentService,
)

func InitApp() *App {

	wire.Build(
		// 第三方依赖
		thirdPartySet,
		commentSvcSet,

		grpc.NewArticleCommentServiceServer,
		ioc.InitGRpcServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
