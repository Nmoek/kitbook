package repository

import (
	"context"
	"kitbook/tag/domain"
	"kitbook/tag/repository/cache"
	"kitbook/tag/repository/dao"
	"time"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag domain.Tag) (int64, error)
	GetTags(ctx context.Context, uid int64) ([]domain.Tag, error)
}

type CacheTagRepository struct {
	dao   dao.TagDao
	cache cache.TagCache
}

func NewCacheTagRepository(dao dao.TagDao, cache cache.TagCache) TagRepository {
	return &CacheTagRepository{
		dao:   dao,
		cache: cache,
	}
}

func (c *CacheTagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	id, err := c.dao.CreateTag(ctx, c.toTagDao(&tag))
	if err != nil {
		return 0, err
	}

	err = c.cache.Append(ctx, id, tag)
	if err != nil {
		//TODO:日志告警
	}
	return id, nil
}

func (c *CacheTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {

	res, err := c.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}

	//TODO: 日志告警

	tags, err := c.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}

	res = make([]domain.Tag, 0, len(tags))
	for _, tag := range tags {
		res = append(res, c.toTagDomain(&tag))
	}

	err = c.cache.Append(ctx, uid, res...)
	if err != nil {
		//TODO:日志告警
	}
	return res, nil
}

// @func: PreloadUserTags
// @date: 2024-03-12 00:58:32
// @brief: ToB业务标签预加载
// @author: Kewin Li
// @receiver c
// @param ctx
// @return error
func (c *CacheTagRepository) PreloadUserTags(ctx context.Context) error {
	//TODO: 在服务初始化时就进行全局预加载, 放到ioc去做
	offset := 0
	batch := 100
	for {
		dbCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		tags, err := c.dao.GetTags(dbCtx, offset, batch)
		cancel()
		if err != nil {
			return err
		}
		if len(tags) <= 0 {
			break
		}
		for _, tag := range tags {
			err = c.cache.Append(ctx, tag.Uid, c.toTagDomain(&tag))
			if err != nil {
				//TODO: 日志告警
				continue
			}
		}

		if len(tags) < batch {
			break
		}
		offset += len(tags)
	}

	return nil
}

func (c *CacheTagRepository) toTagDao(tag *domain.Tag) dao.Tag {
	return dao.Tag{
		Uid:  tag.Uid,
		Name: tag.Name,
	}
}

func (c *CacheTagRepository) toTagDomain(tag *dao.Tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Uid:  tag.Uid,
		Name: tag.Name,
	}
}
