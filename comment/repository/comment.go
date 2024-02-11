package repository

import (
	"context"
	"database/sql"
	"kitbook/comment/domain"
	"kitbook/comment/repository/dao"
	"time"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, cmt domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
}

type ArticleCommentRepository struct {
	dao dao.CommentDao
}

func NewArticleCommentRepository(dao dao.CommentDao) CommentRepository {
	return &ArticleCommentRepository{
		dao: dao,
	}
}

func (a *ArticleCommentRepository) CreateComment(ctx context.Context, cmt domain.Comment) error {
	return a.dao.Insert(ctx, a.ConvertsCommentDao(&cmt))
}
func (a *ArticleCommentRepository) DeleteComment(ctx context.Context, id int64) error {
	return a.dao.Delete(ctx, a.ConvertsCommentDao(&domain.Comment{
		Id: id,
	}))
}

func (a *ArticleCommentRepository) ConvertsCommentDomain(cmt *dao.Comment) domain.Comment {
	res := domain.Comment{
		Id:      cmt.Id,
		Biz:     cmt.Biz,
		BizId:   cmt.BizId,
		Content: cmt.Content,
		Ctime:   time.UnixMilli(cmt.Ctime),
		Utime:   time.UnixMilli(cmt.Utime),
	}

	if cmt.RootId.Valid {
		res.RootComment = &domain.Comment{
			Id: cmt.RootId.Int64,
		}
	}

	if cmt.Pid.Valid {
		res.ParentComment = &domain.Comment{
			Id: cmt.Pid.Int64,
		}

	}
	return res
}

func (a *ArticleCommentRepository) ConvertsCommentDao(cmt *domain.Comment) dao.Comment {
	res := dao.Comment{
		Id:      cmt.Id,
		Biz:     cmt.Biz,
		BizId:   cmt.BizId,
		Content: cmt.Content,

		Utime: time.Now().UnixMilli(),
		Ctime: time.Now().UnixMilli(),
	}

	if cmt.RootComment != nil {
		res.RootId = sql.NullInt64{
			Int64: cmt.RootComment.Id,
			Valid: true,
		}
	}

	if cmt.ParentComment != nil {
		res.Pid = sql.NullInt64{
			Int64: cmt.ParentComment.Id,
			Valid: true,
		}
	}

	return res
}
