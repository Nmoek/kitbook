package grpc

import (
	"context"
	"google.golang.org/grpc"
	rewardv1 "kitbook/api/proto/gen/reward/v1"
	"kitbook/reward/domain"
	"kitbook/reward/service"
)

type RewardServiceServer struct {
	rewardv1.UnimplementedRewardServiceServer
	svc service.RewardService
}

func NewRewardServiceServer(svc service.RewardService) *RewardServiceServer {
	return &RewardServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (r *RewardServiceServer) Register(s *grpc.Server) {
	rewardv1.RegisterRewardServiceServer(s, r)
}

// @func: PreReward
// @date: 2024-02-05 02:42:11
// @brief: xx用户向xx模块的某条记录进行打赏操作
// @author: Kewin Li
// @receiver r
// @param ctx
// @param request
// @return *rewardv1.RewardResponse
// @return error
func (r *RewardServiceServer) PreReward(ctx context.Context, request *rewardv1.PreRewardRequest) (*rewardv1.PreRewardResponse, error) {
	codeURL, err := r.svc.PreReward(ctx, domain.Reward{
		Uid: request.GetUid(),
		Target: domain.Target{
			BizId:   request.GetBizId(),
			Biz:     request.GetBiz(),
			BizName: request.GetBizName(),
			Uid:     request.GetTargetUid(),
		},
		Amt: request.GetAmt(),
	})

	return &rewardv1.PreRewardResponse{
		CodeUrl: codeURL.URL,
		Rid:     codeURL.Rid,
	}, err
}

// @func: UpdateStatus
// @date: 2024-02-05 03:35:13
// @brief: 更新订单支付状态
// @author: Kewin Li
// @receiver r
// @param context.Context
// @param *rewardv1.UpdateStatusRequest
// @return *rewardv1.UpdateStatusResponse
// @return error
func (r *RewardServiceServer) UpdateStatus(ctx context.Context, request *rewardv1.UpdateStatusRequest) (*rewardv1.UpdateStatusResponse, error) {
	err := r.svc.UpdateStatus(ctx, request.GetBizTradeNo(), domain.RewardStatus(request.GetStatus()))
	if err != nil {
		return nil, err
	}

	return &rewardv1.UpdateStatusResponse{}, nil
}

func (r *RewardServiceServer) GetReward(ctx context.Context, request *rewardv1.GetRewardRequest) (*rewardv1.GetRewardResponse, error) {
	_, err := r.svc.GetReward(ctx, request.GetRid(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &rewardv1.GetRewardResponse{}, nil
}
