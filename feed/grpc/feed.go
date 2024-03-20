package grpc

import (
	"google.golang.org/grpc"
	feedv1 "kitbook/api/proto/gen/feed/v1"
	"kitbook/feed/service"
)

type FeedServiceServer struct {
	feedv1.UnimplementedFeedServiceServer
	svc service.FeedService
}

func NewFeedServiceServer(svc service.FeedService) *FeedServiceServer {
	return &FeedServiceServer{
		svc: svc,
	}
}

func (f *FeedServiceServer) Register(server *grpc.Server) {
	feedv1.RegisterFeedServiceServer(server, f)
}
