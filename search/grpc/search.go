package grpc

import (
	"context"
	"google.golang.org/grpc"
	searchv1 "kitbook/api/proto/gen/search/v1"
	"kitbook/search/domain"
	"kitbook/search/service"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	searchSvc service.SearchService
}

func NewSearchServiceServer(searchSvc service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{
		searchSvc: searchSvc,
	}
}

func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	res, err := s.searchSvc.Search(ctx, request.GetUid(), request.GetExpression())
	return &searchv1.SearchResponse{
		User: &searchv1.UserResult{
			Users: s.toRpcUsers(res.Users),
		},
		Article: &searchv1.ArticleResult{
			Articles: s.toRpcArticles(res.Articles),
		},
	}, err
}

func (s *SearchServiceServer) toRpcUsers(users []domain.User) []*searchv1.User {
	res := make([]*searchv1.User, 0, len(users))
	for _, u := range users {
		res = append(res, &searchv1.User{
			Id:       u.Id,
			Email:    u.Email,
			Phone:    u.Phone,
			Nickname: u.Nickname,
		})
	}

	return res
}

func (s *SearchServiceServer) toRpcArticles(arts []domain.Article) []*searchv1.Article {
	res := make([]*searchv1.Article, 0, len(arts))
	for _, a := range arts {

		res = append(res, &searchv1.Article{
			Id:      a.Id,
			Title:   a.Title,
			Content: a.Content,
			Status:  int32(a.Status),
		})
	}

	return res
}
