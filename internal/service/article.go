// Package service
// @Description:
package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/pkg/logger"
)

var ErrInvalidUpdate = errors.New("非法操作")

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx *gin.Context, art domain.Article) error
}

// NormalArticleService
// @Description: 普通帖子服务(免费、免登录)
type NormalArticleService struct {
	//  不分库
	repo repository.ArticleRepository

	/// V1 版本 在service层做数据同步
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository
	///

	l logger.Logger
}

func NewNormalArticleService(repo repository.ArticleRepository, l logger.Logger) ArticleService {
	return &NormalArticleService{
		repo: repo,
	}
}

// @func: NewNormalArticleServiceV1
// @date: 2023-11-25 23:57:48
// @brief: 依赖注入-service层同步数据服务
// @author: Kewin Li
// @param authorRepo
// @param readerRepo
// @param l
// @return ArticleService
func NewNormalArticleServiceV1(
	authorRepo repository.ArticleAuthorRepository,
	readerRepo repository.ArticleReaderRepository,
	l logger.Logger) ArticleService {
	return &NormalArticleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		l:          l,
	}
}

// @func: Save
// @date: 2023-11-23 00:03:50
// @brief: 帖子服务-帖子保存
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
func (n *NormalArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	// 帖子未发表
	art.Status = domain.ArticleStatusUnpublished

	if art.Id > 0 {
		err := n.repo.Update(ctx, art)
		if err == repository.ErrUserMismatch {
			return -1, ErrInvalidUpdate
		}
		return art.Id, err
	}
	return n.repo.Create(ctx, art)
}

// @func: Publish(ctx context.Context, art domain.Article)
// @date: 2023-11-25 23:57:48
// @brief: 帖子服务-帖子发表-dao层数据同步
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
// @return error
func (n *NormalArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	id, err := n.repo.Sync(ctx, art)
	if err == repository.ErrUserMismatch {
		return -1, ErrInvalidUpdate
	}
	return id, err
}

// @func: Publish(ctx context.Context, art domain.Article)
// @date: 2023-11-25 23:57:48
// @brief: 帖子服务-帖子发表-service层同步数据
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
// @return error
func (n *NormalArticleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	// 1. 先存制作库
	// 2. 再存线上库
	var id int64
	var err error
	var logKey = logger.ArticleLogMsgKey[logger.LOG_ART_PUBLISH]
	fields := logger.Fields{}

	if art.Id > 0 {
		err = n.authorRepo.Update(ctx, art)
		if err != nil {
			return art.Id, err
		}
	} else {
		id, err = n.authorRepo.Create(ctx, art)
		if err != nil {
			return 0, err
		}
		art.Id = id
	}

	// insert or uodate 语义
	err = n.readerRepo.Save(ctx, art)
	if err != nil {
		n.l.ERROR(logKey,
			fields.Add(logger.String("制作库保存成功, 线上库保存失败")).
				Add(logger.Error(err)).
				Add(logger.Int[int64]("artId", art.Id)).
				Add(logger.Int[int64]("authId", art.Author.Id))...)
	}
	return art.Id, err
}

// @func: Withdraw
// @date: 2023-11-28 12:49:25
// @brief: 帖子服务-帖子撤回
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
// @return int64
// @return error
func (n *NormalArticleService) Withdraw(ctx *gin.Context, art domain.Article) error {
	return n.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}
