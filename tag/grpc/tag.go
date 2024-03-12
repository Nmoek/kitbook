package grpc

import (
	"context"
	"google.golang.org/grpc"
	tagv1 "kitbook/api/proto/gen/tag/v1"
	"kitbook/tag/domain"
	"kitbook/tag/service"
)

type TagServiceServer struct {
	tagv1.UnimplementedTagServiceServer
	svc service.TagService
}

func NewTagServiceServer(svc service.TagService) *TagServiceServer {
	return &TagServiceServer{
		svc: svc,
	}
}

func (t *TagServiceServer) Register(server *grpc.Server) {
	tagv1.RegisterTagServiceServer(server, t)
}

func (t *TagServiceServer) CreateTag(ctx context.Context, request *tagv1.CreateTagRequest) (*tagv1.CreateTagResponse, error) {
	id, err := t.svc.CreateTag(ctx, request.GetUid(), request.GetName())
	return &tagv1.CreateTagResponse{
		Id: id,
	}, err
}
func (t *TagServiceServer) GetTags(ctx context.Context, request *tagv1.GetTagsRequest) (*tagv1.GetTagsResponse, error) {
	tags, err := t.svc.GetTags(ctx, request.GetUid())
	return &tagv1.GetTagsResponse{
		Tags: t.toTagsRpc(tags),
	}, err
}

// @func: AttachTags
// @date: 2024-03-12 22:13:43
// @brief: 为资源打上标签
// @author: Kewin Li
// @receiver t
// @param ctx
// @param request
// @return *tagv1.AttachTagsResponse
// @return error
func (t *TagServiceServer) AttachTags(ctx context.Context, request *tagv1.AttachTagsRequest) (*tagv1.AttachTagsResponse, error) {
	err := t.svc.AttachTags(ctx, request.GetTids(), request.GetUid(), request.GetBiz(), request.GetBizId())
	return &tagv1.AttachTagsResponse{}, err
}

func (t *TagServiceServer) GetBizTags(ctx context.Context, request *tagv1.GetBizTagsRequest) (*tagv1.GetBizTagsResponse, error) {
	tags, err := t.svc.GetBizTags(ctx, request.GetBizId(), request.GetBiz(), request.GetUid())
	return &tagv1.GetBizTagsResponse{
		Tags: t.toTagsRpc(tags),
	}, err
}

func (t *TagServiceServer) toTagRpc(tag *domain.Tag) *tagv1.Tag {
	return &tagv1.Tag{
		Id:   tag.Id,
		Uid:  tag.Uid,
		Name: tag.Name,
	}
}

func (t *TagServiceServer) toTagsRpc(tags []domain.Tag) []*tagv1.Tag {
	res := make([]*tagv1.Tag, 0, len(tags))
	for _, tag := range tags {
		res = append(res, t.toTagRpc(&tag))
	}

	return res
}
