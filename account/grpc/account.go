package grpc

import (
	"context"
	"google.golang.org/grpc"
	"kitbook/account/domain"
	"kitbook/account/service"
	accountv1 "kitbook/api/proto/gen/account/v1"
)

type AccountServiceServer struct {
	accountv1.UnimplementedAccountServiceServer
	svc service.AccountService
}

func NewAccountServiceServer(svc service.AccountService) *AccountServiceServer {
	return &AccountServiceServer{
		svc: svc,
	}
}

// @func: Register
// @date: 2024-01-09 13:15:39
// @brief: gprc服务注册
// @author: Kewin Li
// @receiver i
// @param s
func (a *AccountServiceServer) Register(s *grpc.Server) {
	accountv1.RegisterAccountServiceServer(s, a)
}

func (a *AccountServiceServer) Credit(ctx context.Context, request *accountv1.CreditRequest) (*accountv1.CreditResponse, error) {
	err := a.svc.Credit(ctx, domain.Credit{
		Biz:   request.GetBiz(),
		BizId: request.GetBizId(),
		Items: a.itemsToDomain(request.GetItems()),
	})
	if err != nil {
		return nil, err
	}

	return &accountv1.CreditResponse{}, nil

}

func (a *AccountServiceServer) itemsToDomain(items []*accountv1.CreditItem) []domain.CreditItem {
	res := make([]domain.CreditItem, 0, len(items))
	for _, item := range items {
		res = append(res, domain.CreditItem{
			Uid:         item.Uid,
			Account:     item.Account,
			AccountType: domain.AccountType(item.AccountType),
			Amt:         item.Amt,
			Currency:    item.Currency,
		})
	}

	return res
}
