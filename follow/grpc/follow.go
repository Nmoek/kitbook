package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	followv1 "kitbook/api/proto/gen/follow/v1"
	"kitbook/follow/service"
)

// NativePaymentServiceServer
// @Description: 评论服务远程调用接口
type FollowServiceServer struct {
	followv1.UnimplementedFollowServiceServer
	svc service.FollowService
}

func NewFollowServiceServer(svc service.FollowService) *FollowServiceServer {
	return &FollowServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-02-12 17:38:11
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver f
// @param s
func (f *FollowServiceServer) Register(s *grpc.Server) {
	followv1.RegisterFollowServiceServer(s, f)
}

func (f *FollowServiceServer) Follow(ctx context.Context, request *followv1.FollowRequest) (*followv1.FollowResponse, error) {
	err := f.svc.Follow(ctx, request.GetFollowee(), request.GetFollower())
	if err != nil {
		return nil, err
	}

	return &followv1.FollowResponse{}, nil
}

func (f *FollowServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	err := f.svc.CancelFollow(ctx, request.GetFollowee(), request.GetFollower())
	if err != nil {
		return nil, err
	}

	return &followv1.CancelFollowResponse{}, nil
}

func (f *FollowServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowInfo not implemented")
}

func (f *FollowServiceServer) GetFollowStatics(ctx context.Context, request *followv1.GetFollowStaticsRequest) (*followv1.GetFollowStaticsResponse, error) {
	statics, err := f.svc.GetFollowStatics(ctx, request.GetUid())
	if err != nil {
		return nil, err
	}

	return &followv1.GetFollowStaticsResponse{
		Followees: statics.Followees,
		Followers: statics.Followers,
	}, nil
}
