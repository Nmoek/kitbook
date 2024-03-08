package grpc

import (
	"context"
	"google.golang.org/grpc"
	commentv1 "kitbook/api/proto/gen/comment/v1"
	"kitbook/comment/domain"
	"kitbook/comment/service"
	"math"
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

// @func: CreateComment
// @date: 2024-02-11 14:14:45
// @brief: 远程调用服务-创建评论
// @author: Kewin Li
// @receiver a
// @param ctx
// @param request
// @return *commentv1.CreateCommentResponse
// @return error
func (a *ArticleCommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := a.svc.CreateComment(ctx, a.toDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

// @func: DeleteComment
// @date: 2024-02-11 14:15:04
// @brief: 远程调用服务-删除评论
// @author: Kewin Li
// @receiver a
// @param ctx
// @param request
// @return *commentv1.DeleteCommentResponse
// @return error
func (a *ArticleCommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := a.svc.DeleteComment(ctx, request.GetId())
	return &commentv1.DeleteCommentResponse{}, err
}

// @func: GetCommentList
// @date: 2024-02-11 14:15:55
// @brief: 远程调用服务-加载评论
// @author: Kewin Li
// @receiver a
// @param ctx
// @param request
// @return *commentv1.GetCommentListResponse
// @return error
func (a *ArticleCommentServiceServer) GetCommentList(ctx context.Context, request *commentv1.GetCommentListRequest) (*commentv1.GetCommentListResponse, error) {
	minId := request.GetMinId()
	// 第一次查询
	if minId <= 0 {
		minId = math.MaxInt64
	}
	cmtsDomain, err := a.svc.GetCommentList(ctx,
		request.GetBizId(),
		request.GetBiz(),
		minId,
		request.GetLimit())
	if err != nil {
		return nil, err
	}

	return &commentv1.GetCommentListResponse{
		Comments: a.toDTOs(cmtsDomain),
	}, nil

}

func (a *ArticleCommentServiceServer) GetMoreReplies(ctx context.Context, request *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	maxId := request.GetMaxId()
	if maxId <= 0 {
		maxId = math.MinInt64
	}
	cmtsDomain, err := a.svc.GetMoreReplies(ctx,
		request.GetRootId(),
		maxId,
		request.GetLimit())
	if err != nil {
		return nil, err
	}

	return &commentv1.GetMoreRepliesResponse{
		Comments: a.toDTOs(cmtsDomain),
	}, nil
}

func (a *ArticleCommentServiceServer) toDomain(cmt *commentv1.Comment) domain.Comment {
	return domain.Comment{}
}

func (a *ArticleCommentServiceServer) toDTO(cmt *domain.Comment) *commentv1.Comment {
	res := &commentv1.Comment{
		Id:      cmt.Id,
		Biz:     cmt.Biz,
		BizId:   cmt.BizId,
		Content: cmt.Content,
		Utime:   cmt.Utime.UnixMilli(),
		Ctime:   cmt.Ctime.UnixMilli(),
	}
	if cmt.RootComment != nil {
		res.RootComment = a.toDTO(cmt.RootComment)
	}

	if cmt.ParentComment != nil {
		res.ParentComment = a.toDTO(cmt.ParentComment)
	}

	return res

}

func (a *ArticleCommentServiceServer) toDTOs(cmts []domain.Comment) []*commentv1.Comment {
	res := make([]*commentv1.Comment, 0, len(cmts))
	for _, cmt := range cmts {
		res = append(res, a.toDTO(&cmt))
	}

	return res
}

func (a *ArticleCommentServiceServer) isHotBiz(bizId int64, biz string) bool {
	//TODO: 通过热榜模块那边来知道是否是热门资源
	return true
}
