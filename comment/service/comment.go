package service

import (
	"context"
	"kitbook/comment/domain"
	"kitbook/comment/repository"
	"kitbook/pkg/logger"
)

type ArticleCommentService struct {
	repo repository.CommentRepository

	l logger.Logger
}

func NewArticlePaymentService(repo repository.CommentRepository, l logger.Logger) CommentService {
	return &ArticleCommentService{
		repo: repo,
		l:    l,
	}
}

func (a *ArticleCommentService) CreateComment(ctx context.Context, cmt domain.Comment) error {
	return a.repo.CreateComment(ctx, cmt)
}

func (a *ArticleCommentService) DeleteComment(ctx context.Context, id int64) error {
	return a.repo.DeleteComment(ctx, id)
}
