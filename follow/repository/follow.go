package repository

import (
	"context"
	"kitbook/follow/domain"
	"kitbook/follow/repository/cache"
	"kitbook/follow/repository/dao"
	"kitbook/pkg/logger"
)

type FollowRepository interface {
	AddFollowRelation(ctx context.Context, f domain.FollowRelation) error
	InactiveFollowRelation(ctx context.Context, followee int64, follower int64) error
	GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error)
	FollowInfo(ctx context.Context, followee, follower int64) (domain.FollowRelation, error)

	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

type CacheFollowRepository struct {
	dao   dao.FollowDao
	cache cache.FollowCache

	l logger.Logger
}

func NewCacheFollowRepository(dao dao.FollowDao, cache cache.FollowCache, l logger.Logger) FollowRepository {
	return &CacheFollowRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (c *CacheFollowRepository) AddFollowRelation(ctx context.Context, f domain.FollowRelation) error {
	err := c.dao.CreateFollowRelation(ctx, c.ConvertsFRDao(&f))
	if err != nil {
		return err
	}

	// 缓存增加的关注数+粉丝数
	return c.cache.Follow(ctx, f.Follower, f.Followee)
}

func (c *CacheFollowRepository) InactiveFollowRelation(ctx context.Context, followee int64, follower int64) error {
	err := c.dao.UpdateStatus(ctx, followee, follower, dao.FollowRelationStatusInactive)
	if err != nil {
		return err
	}

	// 缓存减少的关注数+粉丝数
	return c.cache.CancelFollow(ctx, follower, followee)
}

// func (c *CacheFollowRepository) ConvertsCommentDomain(cmt *dao.Comment) domain.Comment {
//
// }
func (c *CacheFollowRepository) ConvertsFRDao(f *domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{}
}

func (c *CacheFollowRepository) GetFollowee(ctx context.Context, follower int64, offset int64, limit int64) ([]domain.FollowRelation, error) {
	// 需要缓存吗？ 意义不大, 可以考虑只缓存第一页数据

	// 查数据库
	list, err := c.dao.FollowRelationList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}

	return c.genFollowRelationList(list), nil

}

func (c *CacheFollowRepository) FollowInfo(ctx context.Context, followee, follower int64) (domain.FollowRelation, error) {

	f, err := c.dao.FollowRelationDetail(ctx, followee, follower)
	if err != nil {
		return domain.FollowRelation{}, nil
	}
	return c.convertsDomainFR(&f), nil
}

func (c *CacheFollowRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	// 1. 查缓存 -- 快路径
	res, err := c.cache.GetStaticsInfo(ctx, uid)
	if err == nil {
		return res, err
	}

	c.l.WARN("关注粉丝数量缓存未命中",
		logger.Error(err),
		logger.Int[int64]("uid", uid))

	// 2. 查数据库 -- 慢路径
	res.Followers, err = c.dao.CntFollower(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	res.Followees, err = c.dao.CntFollowee(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}

	// 3. 缓存回写
	err = c.cache.SetStaticsInfo(ctx, uid, res)
	if err != nil {
		c.l.ERROR("关注粉丝数量缓存回写失败",
			logger.Error(err),
			logger.Int[int64]("uid", uid),
			logger.Int[int64]("followers", res.Followers),
			logger.Int[int64]("followees", res.Followees),
		)
	}

	return res, nil
}

func (c *CacheFollowRepository) genFollowRelationList(list []dao.FollowRelation) []domain.FollowRelation {
	res := make([]domain.FollowRelation, 0, len(list))
	for _, item := range list {
		res = append(res, c.convertsDomainFR(&item))
	}

	return res
}

func (c *CacheFollowRepository) convertsDomainFR(f *dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Followee: f.Followee,
		Follower: f.Follower,
	}
}
