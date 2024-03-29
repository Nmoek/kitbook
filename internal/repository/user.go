package repository

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	"kitbook/internal/repository/cache"
	"kitbook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateUser  //邮箱已经存在于表中
	ErrUserNotFound  = dao.ErrRecordNotFound //记录未查询到
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	UpdatePersonalInfo(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openid string) (domain.User, error)
}

type CacheUserRepository struct {
	dao dao.UserDao
	c   cache.UserCache
}

func NewCacheUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao: dao,
		c:   cache,
	}
}

// @func: Create
// @date: 2023-10-08 02:44:45
// @brief: 转发模块-数据创建
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param user
func (repo *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Openid: sql.NullString{
			String: u.WechatInfo.Openid,
			Valid:  u.WechatInfo.Openid != "",
		},
		Unionid: sql.NullString{
			String: u.WechatInfo.Unionid,
			Valid:  u.WechatInfo.Unionid != "",
		},
	})

}

// @func: UpdatePersonalInfo
// @date: 2023-10-13 01:40:56
// @brief: 转发模块-数据修改
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param user
// @return error
func (repo *CacheUserRepository) UpdatePersonalInfo(ctx context.Context, user domain.User) error {
	// TODO: 一致性问题？如何解决？先改redis还是先改数据库？
	return repo.dao.UpdateById(ctx, ConvertsDaoUser(&user))
}

// @func: FindByEmail
// @date: 2023-10-09 01:52:27
// @brief: 转发模块-数据查询
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param email
func (repo *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	findUser, err := repo.dao.FindByEmail(ctx, email)
	if err == gorm.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return ConvertsDomainUser(&findUser), nil
}

// @func: FindById
// @date: 2023-10-12 04:05:17
// @brief: 转发模块-数据查询
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param id
// @return error
func (repo *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 查询缓存
	cacheUser, err := repo.c.Get(ctx, id)

	// 查询缓存err为nil有两种情况:
	//	1. key不存在
	//  2. 网络不通、Redis已经崩溃

	//TODO: 优化细分为key不存在再去查询数据库
	switch err {
	case nil:
		return cacheUser, nil
	case cache.ErrKeyNotExist:
		findUser, err := repo.dao.FindByID(ctx, id)

		if err == gorm.ErrRecordNotFound {
			return domain.User{}, ErrUserNotFound
		}

		if err != nil {
			return domain.User{}, err
		}

		// 插入缓存
		err = repo.c.Set(ctx, ConvertsDomainUser(&findUser))
		// 查询缓存err为nil有两种情况:(缓存穿透)
		// 1. 数据格式等其他因素插入失败
		// 2. 网络不通、Redis已经崩溃

		// TODO: 插入缓存错误日志埋点，不一定要返回错误，会让数据库压力增大但不一定崩溃

		return ConvertsDomainUser(&findUser), nil

	default:
		//TODO: 降级写法
		return domain.User{}, err
	}

}

// @func: FindByPhone
// @date: 2023-10-30 00:10:36
// @brief: 转发模块-数据查询按手机号
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param phone
// @return domain.User
// @return error
func (repo *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {

	findUser, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return ConvertsDomainUser(&findUser), nil
}

// @func: FindByWechat
// @date: 2023-11-12 03:08:29
// @brief: 转发模块-数据查询按微信号
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param openid
// @return domain.User
// @return error
func (repo *CacheUserRepository) FindByWechat(ctx context.Context, openid string) (domain.User, error) {
	findUser, err := repo.dao.FindByWechat(ctx, openid)
	if err != nil {
		return domain.User{}, err
	}

	return ConvertsDomainUser(&findUser), nil
}

// @func: convertsDominUser
// @date: 2023-10-09 02:08:11
// @brief: 转化为domin的User结构体
// @author: Kewin Li
// @param user
// @return domain.User
func ConvertsDomainUser(user *dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: time.UnixMilli(user.Birthday),
		AboutMe:  user.AboutMe,
		WechatInfo: domain.WechatInfo{
			Unionid: user.Unionid.String,
			Openid:  user.Openid.String,
		},
	}
}

// @func: convertsDaoUser
// @date: 2023-10-13 01:56:20
// @brief: 转化为dao的User结构体
// @author: Kewin Li
// @receiver repo
// @param user
// @return dao.User
func ConvertsDaoUser(user *domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Unionid: sql.NullString{
			String: user.WechatInfo.Unionid,
			Valid:  user.WechatInfo.Unionid != "",
		},
		Openid: sql.NullString{
			String: user.WechatInfo.Openid,
			Valid:  user.WechatInfo.Openid != "",
		},
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
}
