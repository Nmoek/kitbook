package repository

import (
	"context"
	"database/sql"
	"golang.org/x/sync/errgroup"
	"kitbook/comment/domain"
	"kitbook/comment/repository/dao"
	"time"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, cmt domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
	FindByBiz(ctx context.Context, bizId int64, biz string, minId int64, limit int64) ([]domain.Comment, error)
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

// @func: FindByBiz
// @date: 2024-02-11 14:51:36
// @brief: 资源加载时查询评论
// @author: Kewin Li
// @receiver a
// @param ctx
// @param bizId
// @param biz
// @param minId
// @param limit
// @return []domain.Comment
// @return error
func (a *ArticleCommentRepository) FindByBiz(ctx context.Context, bizId int64, biz string, minId int64, limit int64) ([]domain.Comment, error) {

	// 1.查询一级评论
	cmtsDao, err := a.dao.FindByBiz(ctx, bizId, biz, minId, limit)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Comment, 0, len(cmtsDao))

	//TODO: 触发限流、降级？ 可以不进行二级评论的查询加载

	// 2. 并发查询二级评论
	var eg errgroup.Group
	for _, cmt := range cmtsDao {
		newCmt := cmt // 规避for变量引用机制引起的问题
		cmtDomain := a.ConvertsCommentDomain(&newCmt)
		res = append(res, cmtDomain)

		eg.Go(func() error {
			// 只查三条
			subCmtsDao, err2 := a.dao.FindRepliesByPid(ctx, cmtDomain.Id, 0, 3)
			if err2 != nil {
				return err2
			}

			// 2.1 装载二级评论
			cmtDomain.Children = make([]domain.Comment, 0, len(subCmtsDao))
			for _, sc := range subCmtsDao {
				cmtDomain.Children = append(cmtDomain.Children, a.ConvertsCommentDomain(&sc))
			}

			return nil
		})
	}

	return res, eg.Wait()

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
