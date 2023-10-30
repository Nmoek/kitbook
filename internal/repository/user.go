package repository

import (
	"context"
	"database/sql"
	"fmt"
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

type UserRepository struct {
	dao *dao.UserDao
	c   *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
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
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
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
func (repo *UserRepository) UpdatePersonalInfo(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.convertsDaoUser(&user))
}

// @func: FindByEmail
// @date: 2023-10-09 01:52:27
// @brief: 转发模块-数据查询
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param email
func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	findUser, err := repo.dao.FindByEmail(ctx, email)
	if err == gorm.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return repo.convertsDomainUser(&findUser), nil
}

// @func: FindByID
// @date: 2023-10-12 04:05:17
// @brief: 转发模块-数据查询
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param id
// @return error
func (repo *UserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	// 查询缓存
	cacheUser, err := repo.c.Get(ctx, id)

	// 查询缓存err为nil有两种情况:
	//	1. key不存在
	//  2. 网络不通、Redis已经崩溃
	fmt.Printf("err:%s, cache user:%v\n", err, cacheUser)

	//TODO: 优化细分为key不存在再去查询数据库
	switch err {
	case nil:
		return repo.convertsDomainUser(&cacheUser), nil
	case cache.ErrKeyNotExist:
		findUser, err := repo.dao.FindByID(ctx, id)

		if err == gorm.ErrRecordNotFound {
			return domain.User{}, ErrUserNotFound
		}

		if err != nil {
			return domain.User{}, err
		}

		// 插入缓存
		err = repo.c.Set(ctx, findUser)
		// 查询缓存err为nil有两种情况:(缓存穿透)
		// 1. 数据格式等其他因素插入失败
		// 2. 网络不通、Redis已经崩溃

		// TODO: 插入缓存错误日志埋点，不一定要返回错误，会让数据库压力增大但不一定崩溃
		fmt.Printf("set key err:%s \n", err)

		return repo.convertsDomainUser(&findUser), nil

	default:
		//TODO: 降级写法
		return domain.User{}, nil
	}

}

// @func: FindByPhone
// @date: 2023-10-30 00:10:36
// @brief: 转发模块-数据查询
// @author: Kewin Li
// @receiver repo
// @param ctx
// @param phone
// @return domain.User
// @return error
func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {

	findUser, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return repo.convertsDomainUser(&findUser), nil
}

// @func: convertsDominUser
// @date: 2023-10-09 02:08:11
// @brief: 转化为domin的User结构体
// @author: Kewin Li
// @param user
// @return domain.User
func (repo *UserRepository) convertsDomainUser(user *dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: time.UnixMilli(user.Birthday),
		AboutMe:  user.AboutMe,
	}
}

// @func: convertsDaoUser
// @date: 2023-10-13 01:56:20
// @brief: 转化为dao的User结构体
// @author: Kewin Li
// @receiver repo
// @param user
// @return dao.User
func (repo *UserRepository) convertsDaoUser(user *domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    sql.NullString{String: user.Email},
		Phone:    sql.NullString{String: user.Phone},
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
}
