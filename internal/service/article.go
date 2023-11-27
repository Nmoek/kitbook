// Package service
// @Description:
package service

import (
	"context"
	"errors"
	"kitbook/internal/domain"
	"kitbook/internal/repository"
	"kitbook/pkg/logger"
)

var ErrInvalidUpdate = errors.New("非法操作")

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

// NormalArticleService
// @Description: 普通帖子服务(免费、免登录)
type NormalArticleService struct {
	// V0 版本 不分库
	repo repository.ArticleRepository

	// V1 版本 在service层做数据同步
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository

	l logger.Logger
}

func NewNormalArticleService(repo repository.ArticleRepository, l logger.Logger) ArticleService {
	return &NormalArticleService{
		repo: repo,
	}
}

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
	if art.Id > 0 {
		err := n.repo.Update(ctx, art)
		if err == repository.ErrInvalidUpdate {
			return -1, ErrInvalidUpdate
		}
		return art.Id, err
	}
	return n.repo.Create(ctx, art)
}

// @func: Publish(ctx context.Context, art domain.Article)
// @date: 2023-11-25 23:57:48
// @brief: 帖子服务-帖子发表
// @author: Kewin Li
// @receiver n
// @param ctx
// @param art
// @return error
func (n *NormalArticleService) PublishV0(ctx context.Context, art domain.Article) (int64, error) {
	return n.repo.Sync(ctx, art)
}

func (n *NormalArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
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
