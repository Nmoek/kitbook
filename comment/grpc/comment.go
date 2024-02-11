package grpc

import (
	"context"
	"google.golang.org/grpc"
	commentv1 "kitbook/api/proto/gen/comment/v1"
	"kitbook/comment/domain"
	"kitbook/comment/service"
)

// NativePaymentServiceServer
// @Description: 评论服务远程调用接口
type ArticleCommentServiceServer struct {
	commentv1.CommentServiceServer
	svc service.CommentService
}

func NewArticleCommentServiceServer(svc service.CommentService) *ArticleCommentServiceServer {
	return &ArticleCommentServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (a *ArticleCommentServiceServer) Register(s *grpc.Server) {
	commentv1.RegisterCommentServiceServer(s, a)
}

func (a *ArticleCommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := a.svc.CreateComment(ctx, a.toDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}
func (a *ArticleCommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := a.svc.DeleteComment(ctx, request.GetId())
	return &commentv1.DeleteCommentResponse{}, err
}

func (a *ArticleCommentServiceServer) toDomain(cmt *commentv1.Comment) domain.Comment {
	return domain.Comment{}
}
