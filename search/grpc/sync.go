package grpc

import (
	"context"
	"google.golang.org/grpc"
	searchv1 "kitbook/api/proto/gen/search/v1"
	"kitbook/search/domain"
	"kitbook/search/service"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	syncSvc service.SyncService
}

func NewSyncServiceServer(syncSvc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{
		syncSvc: syncSvc,
	}
}

func (s *SyncServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func (s *SyncServiceServer) InputUser(ctx context.Context, request *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	err := s.syncSvc.InputUser(ctx, s.toDomainUser(request.GetUser()))
	return &searchv1.InputUserResponse{}, err
}
func (s *SyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.syncSvc.InputArticle(ctx, s.toDomainArticle(request.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) InputAny(ctx context.Context, request *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.syncSvc.Input(ctx, request.GetIndexName(), request.GetDocId(), request.GetData())
	return &searchv1.InputAnyResponse{}, err
}

func (s *SyncServiceServer) toDomainUser(user *searchv1.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
	}
}

func (s *SyncServiceServer) toDomainArticle(art *searchv1.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  int8(art.Status),
	}
}
