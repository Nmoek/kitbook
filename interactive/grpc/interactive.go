package grpc

import (
	"context"
	"google.golang.org/grpc"
	"kitbook/api/proto/gen/intr/v1"
	"kitbook/interactive/domain"
	"kitbook/interactive/service"
)

type InteractiveServiceServer struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (i *InteractiveServiceServer) Register(s *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(s, i)
}

// @func: IncreaseReadCnt
// @date: 2024-01-09 01:29:52
// @brief: 远程执行-阅读数+1
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.IncreaseReadCntResponse
// @return error
func (i *InteractiveServiceServer) IncreaseReadCnt(ctx context.Context, request *intrv1.IncreaseReadCntRequest) (*intrv1.IncreaseReadCntResponse, error) {
	err := i.svc.IncreaseReadCnt(ctx, request.GetBiz(), request.GetBizId())
	if err != nil {
		return nil, err
	}
	return &intrv1.IncreaseReadCntResponse{}, nil
}

// @func: Like
// @date: 2024-01-09 01:30:38
// @brief: 远程执行-点赞数+1
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.LikeResponse
// @return error
func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &intrv1.LikeResponse{}, nil
}

// @func: CancelLike
// @date: 2024-01-09 01:30:51
// @brief: 远程执行-点赞数-1
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.CancelLikeResponse
// @return error
func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &intrv1.CancelLikeResponse{}, nil
}

// @func: Collect
// @date: 2024-01-09 01:31:06
// @brief: 远程执行-收藏数+1
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.CollectResponse
// @return error
func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.GetBiz(), request.GetBizId(), request.GetCollectId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &intrv1.CollectResponse{}, nil
}

// @func: CancelCollect
// @date: 2024-01-09 01:31:38
// @brief: 远程执行-收藏数-1
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.CancelCollectResponse
// @return error
func (i *InteractiveServiceServer) CancelCollect(ctx context.Context, request *intrv1.CancelCollectRequest) (*intrv1.CancelCollectResponse, error) {
	err := i.svc.CancelCollect(ctx, request.GetBiz(), request.GetBizId(), request.GetCollectId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &intrv1.CancelCollectResponse{}, nil
}

// @func: Get
// @date: 2024-01-09 01:32:09
// @brief: 远程执行-聚合查询
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.GetResponse
// @return error
func (i *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	intr, err := i.svc.Get(ctx, request.GetBiz(), request.GetBizId(), request.GetUserId())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.ConvertsAPIInteractive(&intr),
	}, nil
}

// @func: GetByIds
// @date: 2024-01-09 01:38:23
// @brief: 远程执行-批量查询
// @author: Kewin Li
// @receiver i
// @param ctx
// @param request
// @return *intrv1.GetByIdsResponse
// @return error
func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetBizIds())
	if err != nil {
		return nil, err
	}

	intrs := make(map[int64]*intrv1.Interactive, len(res))
	for _, intr := range res {
		intrs[intr.BizId] = i.ConvertsAPIInteractive(&intr)
	}

	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, nil
}

// @func: ConvertsAPIInteractive
// @date: 2024-01-09 01:36:01
// @brief: domain->api 互动表数据
// @author: Kewin Li
// @receiver i
// @param intr
// @return intrv1.Interactive
func (i *InteractiveServiceServer) ConvertsAPIInteractive(intr *domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		BizId:      intr.BizId,
		Biz:        intr.Biz,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
	}
}
