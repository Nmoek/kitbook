package grpc

import (
	"context"
	"errors"
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

// @func: GetCommentList
// @date: 2024-02-11 16:17:09
// @brief: 针对热门资源情况进行不同限流
// @author: Kewin Li
// @receiver a
// @param ctx
// @param request
// @return *commentv1.GetCommentListResponse
// @return error
func (a *AsyncCommentService) GetCommentList(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {

	isHotBiz := a.isHotBiz(request.GetBizId(), request.GetBiz())
	// 使用不同的限流器
	if isHotBiz {
		// 400/s 限流
	} else {
		// 100/s 限流
	}

	return a.ArticleCommentServiceServer.GetCommentList(ctx, request)
}

// @func: GetCommentListV1
// @date: 2024-02-11 16:18:22
// @brief: 限流+降级
// @author: Kewin Li
// @receiver a
// @param ctx
// @param request
// @return *commentv1.GetCommentListResponse
// @return error
func (a *AsyncCommentService) GetCommentListV1(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {

	isHotBiz := a.isHotBiz(request.GetBizId(), request.GetBiz())
	// 非热门数据访问 + 降级状态 不进行数据查询
	if !isHotBiz && ctx.Value("downgrade") == "true" {
		return nil, errors.New("非热门数据处于降级状态")

	}

	return a.ArticleCommentServiceServer.GetCommentList(ctx, request)
}
