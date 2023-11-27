package repository

import (
	"context"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
	"kitbook/pkg/logger"
)

var ErrInvalidUpdate = dao.ErrInvalidUpdate

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type CacheArticleRepository struct {
	// V0写法 不分库
	dao dao.ArticleDao

	// V2写法 在repository层做数据同步
	authorDao dao.ArticleAuthorDao
	readerDao dao.ArticleReaderDao
	db        *gorm.DB
	l         logger.Logger
}

func NewCacheArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}

func NewCacheArticleRepositoryV2(
	authorDao dao.ArticleAuthorDao,
	readerDao dao.ArticleReaderDao,
	db *gorm.DB,
	l logger.Logger) *CacheArticleRepository {
	return &CacheArticleRepository{
		authorDao: authorDao,
		readerDao: readerDao,
		db:        db,
		l:         l,
	}
}

// @func: Create
// @date: 2023-11-24 21:01:34
// @brief: 新建帖子并保存
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, ConvertsDaoArticle(&art))
}

// @func: Update
// @date: 2023-11-24 21:01:45
// @brief: 修改帖子
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return error
func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, ConvertsDaoArticle(&art))
}

// @func: Sync(ctx context.Context, art domain.Article)
// @date: 2023-11-26 19:13:56
// @brief: 帖子发表-数据同步-非事务实现
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {

	var id int64
	var err error
	artDao := ConvertsDaoArticle(&art)

	if art.Id > 0 {
		err = c.authorDao.Update(ctx, artDao)
		if err != nil {
			return art.Id, err
		}

	} else {

		id, err = c.authorDao.Create(ctx, artDao)
		if err != nil {
			return 0, err
		}
		art.Id = id
	}

	err = c.readerDao.Upsert(ctx, artDao)
	return art.Id, err
}

// @func: Sync(ctx context.Context, art domain.Article)
// @date: 2023-11-26 19:13:56
// @brief: 帖子发表-数据同步-非事务实现
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {

	var id int64
	var err error
	artDao := ConvertsDaoArticle(&art)

	// 开启事务
	dbCtx := c.db.WithContext(ctx).Begin()
	if dbCtx.Error != nil {
		return 0, dbCtx.Error
	}
	defer dbCtx.Rollback() //事务回滚

	authorDao := dao.NewGormArticleAuthorDao(dbCtx)
	readerDao := dao.NewGormArticleReaderDao(dbCtx)

	if art.Id > 0 {
		err = authorDao.Update(ctx, artDao)
		if err != nil {
			return art.Id, err
		}

	} else {

		id, err = authorDao.Create(ctx, artDao)
		if err != nil {
			return 0, err
		}
		art.Id = id
	}

	// 操作另外一张表
	err = readerDao.UpsertV2(ctx, dao.PublishedArticle(artDao))

	err = dbCtx.Commit().Error
	return art.Id, err

}

// @func: Sync
// @date: 2023-11-27 12:45:29
// @brief: 数据转发-数据同步-dao层同步
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, ConvertsDaoArticle(&art))
}

// @func: convertsDominUser
// @date: 2023-10-09 02:08:11
// @brief: 转化为domin的Article结构体
// @author: Kewin Li
// @param user
// @return domain.User
func ConvertsDomainArticle(art *dao.Article) domain.Article {
	return domain.Article{}
}

// @func: ConvertsDaoUser
// @date: 2023-11-23 00:55:51
// @brief: 转化为dao的Article结构体
// @author: Kewin Li
// @param user
// @return dao.User
func ConvertsDaoArticle(art *domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
