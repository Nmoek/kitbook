package grpc

import (
	"context"
	commentv1 "kitbook/api/proto/gen/comment/v1"
	"kitbook/comment/events"
)

type AsyncCommentService struct {
	ArticleCommentServiceServer
	producer events.Producer
}

func NewAsyncCommentService(producer events.Producer) *AsyncCommentService {
	return &AsyncCommentService{
		producer: producer,
	}
}

func (a *AsyncCommentService) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	// 触发限流或降级进行kafka消息发送
	if ctx.Value("limited") == "true" || ctx.Value("downgrade") == "true" {

	}

	err := a.svc.CreateComment(ctx, a.toDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}
func (a *AsyncCommentService) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	// 触发限流或降级进行kafka消息发送
	if ctx.Value("limited") == "true" || ctx.Value("downgrade") == "true" {

	}
	err := a.svc.DeleteComment(ctx, request.GetId())
	return &commentv1.DeleteCommentResponse{}, err
}
