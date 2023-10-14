package repository

import (
	"context"
	"gorm.io/gorm"
	"kitbook/internal/domain"
	"kitbook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail //邮箱已经存在于表中
	ErrUserNotFound   = dao.ErrRecordNotFound //记录未查询到
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
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
		Email:    u.Email,
		Password: u.Password,
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
	findUser, err := repo.dao.FindByID(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}

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
		Email:    user.Email,
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
		Email:    user.Email,
		Password: user.Password,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
}
