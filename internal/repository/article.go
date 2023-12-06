package repository

import (
	"context"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"kitbook/pkg/logger"
	"time"
)

var ErrUserMismatch = dao.ErrUserMismatch

// 预加载缓存大小限制
const contentLimitSize = 1 * 1024 * 1024

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, artId int64, authorId int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, artId int64) (domain.Article, error)
	GetPubById(ctx context.Context, artId int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	// V0写法 不分库
	dao      dao.ArticleDao
	cache    cache.ArticleCache
	userRepo UserRepository

	// V2写法 在repository层做数据同步
	authorDao dao.ArticleAuthorDao
	readerDao dao.ArticleReaderDao
	db        *gorm.DB
	l         logger.Logger
}

func NewCacheArticleRepository(dao dao.ArticleDao,
	cache cache.ArticleCache,
	userRepo UserRepository) ArticleRepository {
	return &CacheArticleRepository{
		dao:      dao,
		cache:    cache,
		userRepo: userRepo,
	}
}

// @func: NewCacheArticleRepositoryV2
// @date: 2023-12-02 18:28:36
// @brief: 依赖注入-repository层同步数据服务
// @author: Kewin Li
// @param authorDao
// @param readerDao
// @param db
// @param l
// @return *CacheArticleRepository
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
	err := c.dao.UpdateById(ctx, ConvertsDaoArticle(&art))
	if err == nil {
		err = c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			//TODO: 日志埋点
		}
	}
	//TODO: 日志埋点

	return err
}

// @func: Sync(ctx context.Context, art domain.Article)
// @date: 2023-11-26 19:13:56
// @brief: 帖子发表-数据同步-repository非事务实现(不同库)
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
// @brief: 帖子发表-数据同步-repository事务实现(同库不同表)
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
// @brief: 帖子发表-数据同步-dao层同步
// @author: Kewin Li
// @receiver c
// @param ctx
// @param art
// @return int64
// @return error
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, ConvertsDaoArticle(&art))
	if err == nil {
		err = c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			//TODO: 日志埋点
		}
	}

	// 进行发表时缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// TODO: 发表缓存优化点
		// 1. 根据不同创作者的火热程度设置过期时间
		// 2. 对新用户进行帖子流量倾斜来设置过期时间
		// 3. 对内容的优质程度来设置过期时间
		// ... 等等 这些角度的考虑都是对内容搜推有着极大影响！！！
		author, err := c.userRepo.FindById(ctx, art.Author.Id)
		if err != nil {
			// TODO: 日志埋点
			return
		}
		art.Author = domain.Author{
			Id:   author.Id,
			Name: author.Nickname,
		}

		err = c.cache.SetPub(ctx, art)
		if err != nil {
			// TODO: 日志埋点
		}

	}()

	//TODO: 日志埋点
	return id, err
}

// @func: SyncStatus
// @date: 2023-11-28 12:52:50
// @brief: 帖子撤回-状态同步-dao层同步
// @author: Kewin Li
// @receiver c
// @param ctx
// @param artId
// @param authorId
// @param status
// @return int64
// @return error
func (c *CacheArticleRepository) SyncStatus(ctx context.Context, artId int64, authorId int64, status domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, artId, authorId, status.ToUint8())
	if err == nil {
		err = c.cache.DelFirstPage(ctx, authorId)
		if err != nil {
			//TODO: 日志埋点
		}
	}

	// TODO: 日志埋点

	return err
}

// @func: GetByAuthor
// @date: 2023-12-04 00:25:08
// @brief: 帖子查询-查询创作者创作列表
// @author: Kewin Li
// @receiver c
// @param ctx
// @param id
// @param offset
// @param limit
// @return []domain.Article
// @return error
func (c *CacheArticleRepository) GetByAuthor(ctx context.Context, userId int64, offset int, limit int) ([]domain.Article, error) {
	// 判定是否应该查询缓存

	// TODO: 优化 只要数据量小于等于1 page都可以进行查询
	if offset == 0 && limit == 100 {
		arts, err := c.cache.GetFirstPage(ctx, userId)
		if err == nil {
			return arts, err
		} else {
			// TODO: 日志埋点, 缓存命中、缓存未命中等情况都需要考虑
		}
	}

	artsDao, err := c.dao.GetByAuthor(ctx, userId, offset, limit)
	arts := make([]domain.Article, len(artsDao))
	for i, art := range artsDao {
		arts[i] = ConvertsDomainArticleFromProduce(&art)
	}

	// 查完数据库需要把缓存放回去
	// TODO:优化  1.异步缓存 2.达到查询阈值才缓存
	if offset == 0 && limit == 100 {
		err = c.cache.SetFirstPage(ctx, userId, arts)
		if err != nil {
			//TODO: 日志埋点
			// 1. 偶尔网络波动缓存失败
			// 2. redis崩溃、网络服务长时间不恢复
			// 3. 操作redis出错
		}
	}

	// 设置缓存预加载
	go func() {
		// 异步之后需要剥离原有的context
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := c.preCache(ctx, arts)
		if err != nil {
			//TODO: 设置缓存预加载出错 日志埋点
		}
	}()

	return arts, err
}

// @func: GetById
// @date: 2023-12-05 02:35:14
// @brief: 帖子查询-查询创作列表详情
// @author: Kewin Li
// @receiver c
// @param ctx
// @param artId
// @return domain.Article
// @return error
func (c *CacheArticleRepository) GetById(ctx context.Context, artId int64) (domain.Article, error) {

	// 检查缓存预加载是否存在
	art, err := c.cache.GetById(ctx, artId)
	if err == nil {
		return art, nil
	}
	// TODO: 检查预加载出错 日志埋点

	artDAO, err := c.dao.GetById(ctx, artId)

	// 库查询结束后 预加载回写
	// 可以同步 也可以异步
	go func() {
		err = c.cache.SetById(ctx, ConvertsDomainArticleFromProduce(&artDAO))
		if err != nil {
			// TODO: 回写预加载出错 日志埋点
		}

	}()

	return ConvertsDomainArticleFromProduce(&artDAO), err
}

// @func: GetPubById
// @date: 2023-12-06 13:14:41
// @brief: 帖子查询-读者查询接口
// @author: Kewin Li
// @receiver c
// @param ctx
// @param artId
// @return domain.Article
// @return error
func (c *CacheArticleRepository) GetPubById(ctx context.Context, artId int64) (domain.Article, error) {
	// 取帖子缓存
	art, err := c.cache.GetPubById(ctx, artId)
	if err == nil {
		return art, err
	}

	// TODO: 取缓存出错、未取到 日志埋点

	res, err := c.dao.GetPubById(ctx, artId)
	if err != nil {
		return domain.Article{}, err
	}

	art = ConvertsDomainArticleFromLive(&res)

	// 帖子缓存回写
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := c.cache.SetPubById(ctx, art)
		if err != nil {
			//TODO: 缓存回写失败 日志埋点
		}

	}()

	//TODO: 查询User信息, 拿到创作者名字
	author, err := c.userRepo.FindById(ctx, art.Author.Id)
	if err != nil {

		return domain.Article{}, err
	}

	art.Author.Name = author.Nickname

	return art, err
}

// @func: preCache
// @date: 2023-12-05 23:13:52
// @brief: 帖子查询-设置列表详情缓存预加载
// @author: Kewin Li
// @receiver c
// @param ctx
// @param arts
// @return interface{}
func (c *CacheArticleRepository) preCache(ctx context.Context, arts []domain.Article) error {
	// 谨慎缓存大文档
	if len(arts) > 0 && len(arts[0].Content) <= contentLimitSize {
		err := c.cache.SetById(ctx, arts[0])
		if err != nil {
			// TODO: 预加载缓存失败 日志埋点
			return err
		}
	}

	return nil
}

// @func: convertsDominUser
// @date: 2023-10-09 02:08:11
// @brief: 制作库转化为domin的Article结构体
// @author: Kewin Li
// @param user
// @return domain.User
func ConvertsDomainArticleFromProduce(art *dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Status: domain.ToArticleStatus(art.Status),
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
	}
}

// @func: ConvertsDomainArticleFromLive
// @date: 2023-12-07 02:10:00
// @brief: 线上库转化为domain的Article结构体
// @author: Kewin Li
// @param art
// @return domain.Article
func ConvertsDomainArticleFromLive(art *dao.PublishedArticle) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
			// 深度理解这里为什么不默认进行创作者名称赋值
		},
		Status: domain.ToArticleStatus(art.Status),
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
	}
}

// @func: ConvertsDaoUser
// @date: 2023-11-23 00:55:51
// @brief: 转化为dao的制作库Article结构体
// @author: Kewin Li
// @param user
// @return dao.User
func ConvertsDaoArticle(art *domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

// @func: ConvertsDaoPublishedArticle
// @date: 2023-12-06 23:11:37
// @brief: 转化为dao的线上库Article结构体
// @author: Kewin Li
// @param art
// @return dao.PublishedArticle
func ConvertsDaoPublishedArticle(art *domain.Article) dao.PublishedArticle {
	return dao.PublishedArticle{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
