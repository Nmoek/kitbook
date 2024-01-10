package client

import (
	"context"
	"google.golang.org/grpc"
	intrv1 "kitbook/api/proto/gen/intr/v1"
	"kitbook/interactive/domain"
	"kitbook/interactive/service"
)

// LocalInteractiveServiceAdapter
// @Description:适配器模式-实际本地调用伪装成远程调用
type LocalInteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewLocalInteractiveServiceAdapter(svc service.InteractiveService) *LocalInteractiveServiceAdapter {
	return &LocalInteractiveServiceAdapter{
		svc: svc,
	}
}

func (l *LocalInteractiveServiceAdapter) IncreaseReadCnt(ctx context.Context, in *intrv1.IncreaseReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncreaseReadCntResponse, error) {
	err := l.svc.IncreaseReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncreaseReadCntResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := l.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUserId())
	return &intrv1.LikeResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := l.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUserId())
	return &intrv1.CancelLikeResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := l.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetCollectId(), in.GetUserId())
	return &intrv1.CollectResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) CancelCollect(ctx context.Context, in *intrv1.CancelCollectRequest, opts ...grpc.CallOption) (*intrv1.CancelCollectResponse, error) {
	err := l.svc.CancelCollect(ctx, in.GetBiz(), in.GetBizId(), in.GetCollectId(), in.GetUserId())
	return &intrv1.CancelCollectResponse{}, err
}

func (l *LocalInteractiveServiceAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	intr, err := l.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUserId())
	return &intrv1.GetResponse{
		Intr: l.ConvertsAPIInteractive(&intr),
	}, err
}

func (l *LocalInteractiveServiceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	res, err := l.svc.GetByIds(ctx, in.GetBiz(), in.GetBizIds())

	intrs := map[int64]*intrv1.Interactive{}

	for _, intr := range res {
		intrs[intr.BizId] = l.ConvertsAPIInteractive(&intr)
	}

	return &intrv1.GetByIdsResponse{
		Intrs: intrs,
	}, err
}

func (l *LocalInteractiveServiceAdapter) ConvertsAPIInteractive(intr *domain.Interactive) *intrv1.Interactive {
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
