package service

import (
	"context"
	"kitbook/comment/domain"
)

type CommentService interface {
	CreateComment(ctx context.Context, cmt domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	GetCommentList(ctx context.Context, bizId int64, biz string, minId int64, limit int64) ([]domain.Comment, error)
	GetMoreReplies(ctx context.Context, rootId int64, maxId int64, limit int64) ([]domain.Comment, error)
}
