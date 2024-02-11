package service

import (
	"context"
	"kitbook/comment/domain"
)

type CommentService interface {
	CreateComment(ctx context.Context, cmt domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
}
