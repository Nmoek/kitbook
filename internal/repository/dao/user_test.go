// Package dao
// @Description: 单元测试-用户模块
package dao

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	msqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

var tableColunms = []string{
	"id",
	"email",
	"phone",
	"openid",
	"unionid",
	"password",
	"nickname",
	"birthday",
	"about_me",
	"ctime",
	"utime",
}

type GormUserDaoSuite struct {
	suite.Suite
}

// @func: TestInsert
// @date: 2024-01-25 03:55:07
// @brief: 单元测试-dao层-插入用户数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestInsert() {
	t := g.T()

	now := time.Now().UnixMilli()

	testCases := []struct {
		name string

		mock func(t *testing.T) *sql.DB

		user User

		wantErr error
	}{
		// 数据插入成功
		{
			name: "Insert data successfully",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(mockRes)

				return db
			},

			user: User{
				Email:    sql.NullString{String: "1@qq.com", Valid: true},
				Phone:    sql.NullString{String: "111111111", Valid: true},
				Openid:   sql.NullString{String: "1", Valid: true},
				Unionid:  sql.NullString{String: "1", Valid: true},
				Nickname: "test",

				Password: "**********",
				Birthday: now,
				AboutMe:  "test content",
				Ctime:    now,
				Utime:    now,
			},
		},
		// 唯一索引冲突
		{
			name: "Unique Index Conflict",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(&msqlDriver.MySQLError{Number: 1062})

				return db
			},

			user: User{
				Email:    sql.NullString{String: "2@qq.com", Valid: true},
				Phone:    sql.NullString{String: "222222222", Valid: true},
				Openid:   sql.NullString{String: "2", Valid: true},
				Unionid:  sql.NullString{String: "2", Valid: true},
				Nickname: "test",

				Password: "**********",
				Birthday: now,
				AboutMe:  "test content",
				Ctime:    now,
				Utime:    now,
			},

			wantErr: ErrDuplicateUser,
		},
		// 其他错误
		{
			name: "Unique Index Conflict",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
						sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("数据库错误"))

				return db
			},

			user: User{

				Nickname: "test",
			},

			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true, //是否跳过版本查询
			}), &gorm.Config{
				DisableAutomaticPing:   true, // 是否禁止ping数据库
				SkipDefaultTransaction: true, // 是否禁止事务
			})
			assert.NoError(t, err)

			d := NewGormUserDao(db)
			err = d.Insert(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// @func: TestFindByEmail
// @date: 2024-01-25 03:55:49
// @brief: 单元测试-dao层-按邮箱查询数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestFindByEmail() {
	t := g.T()
	//now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(t *testing.T, u User) *sql.DB

		email    string
		user     User
		wantUser User
		wantErr  error
	}{
		// 查询成功
		{
			name: "Query successfully",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("1@qq.com").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(nil)

				return db

			},
			email: "1@qq.com",
			user: User{
				Id:    1,
				Email: sql.NullString{String: "1@qq.com", Valid: true},
			},
			wantUser: User{
				Id:    1,
				Email: sql.NullString{String: "1@qq.com", Valid: true},
			},
		},
		// 没有记录
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("2@qq.com").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(gorm.ErrRecordNotFound)

				return db

			},
			email:   "2@qq.com",
			wantErr: gorm.ErrRecordNotFound,
		},
		// 其他错误
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("3@qq.com").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(errors.New("数据库错误"))

				return db

			},
			email:   "3@qq.com",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t, tc.user)
			defer sqlDB.Close()

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})

			assert.NoError(t, err)
			d := NewGormUserDao(db)
			user, err := d.FindByEmail(context.Background(), tc.email)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})

	}
}

// @func: TestFindByID
// @date: 2024-01-25 21:50:58
// @brief:  单元测试-dao层-按用户id查询数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestFindByID() {
	t := g.T()
	//now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(t *testing.T, u User) *sql.DB

		id       int64
		user     User
		wantUser User
		wantErr  error
	}{
		// 查询成功
		{
			name: "Query successfully",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(nil)

				return db

			},
			id: 1,
			user: User{
				Id:    1,
				Email: sql.NullString{String: "1@qq.com", Valid: true},
			},
			wantUser: User{
				Id:    1,
				Email: sql.NullString{String: "1@qq.com", Valid: true},
			},
		},
		// 没有记录
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(gorm.ErrRecordNotFound)

				return db

			},
			id:      2,
			wantErr: gorm.ErrRecordNotFound,
		},
		// 其他错误
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs(3).
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(errors.New("数据库错误"))

				return db

			},
			id:      3,
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t, tc.user)
			defer sqlDB.Close()

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})

			assert.NoError(t, err)
			d := NewGormUserDao(db)
			user, err := d.FindByID(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})

	}
}

// @func: TestFindByPhone
// @date: 2024-01-25 21:58:59
// @brief: 单元测试-dao层-按手机号查询数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestFindByPhone() {
	t := g.T()
	//now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(t *testing.T, u User) *sql.DB

		phone    string
		user     User
		wantUser User
		wantErr  error
	}{
		// 查询成功
		{
			name: "Query successfully",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("111").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(nil)

				return db

			},
			phone: "111",
			user: User{
				Id:    1,
				Phone: sql.NullString{String: "111", Valid: true},
			},
			wantUser: User{
				Id:    1,
				Phone: sql.NullString{String: "111", Valid: true},
			},
		},
		// 没有记录
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("222").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(gorm.ErrRecordNotFound)

				return db

			},
			phone:   "222",
			wantErr: gorm.ErrRecordNotFound,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("333").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(errors.New("数据库错误"))

				return db

			},
			phone:   "333",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t, tc.user)
			defer sqlDB.Close()

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})

			assert.NoError(t, err)
			d := NewGormUserDao(db)
			user, err := d.FindByPhone(context.Background(), tc.phone)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})

	}
}

// @func: TestFindByWechat
// @date: 2024-01-25 23:02:32
// @brief: 单元测试-dao层-按微信号查询数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestFindByWechat() {
	t := g.T()
	//now := time.Now().UnixMilli()
	testCases := []struct {
		name string

		mock func(t *testing.T, u User) *sql.DB

		openid   string
		user     User
		wantUser User
		wantErr  error
	}{
		// 查询成功
		{
			name: "Query successfully",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("111").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(nil)

				return db

			},
			openid: "111",
			user: User{
				Id:     1,
				Openid: sql.NullString{String: "111", Valid: true},
			},
			wantUser: User{
				Id:     1,
				Openid: sql.NullString{String: "111", Valid: true},
			},
		},
		// 没有记录
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("222").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(gorm.ErrRecordNotFound)

				return db

			},
			openid:  "222",
			wantErr: gorm.ErrRecordNotFound,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT .*").
					WithArgs("333").
					WillReturnRows(sqlmock.NewRows(tableColunms).AddRow(
						u.Id,
						u.Email, u.Phone, u.Openid, u.Unionid,
						u.Password, u.Nickname, u.Birthday, u.AboutMe, u.Ctime, u.Utime,
					)).WillReturnError(errors.New("数据库错误"))

				return db

			},
			openid:  "333",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t, tc.user)
			defer sqlDB.Close()

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})

			assert.NoError(t, err)
			d := NewGormUserDao(db)
			user, err := d.FindByWechat(context.Background(), tc.openid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})

	}
}

// @func: TestUpdateById
// @date: 2024-01-25 22:48:45
// @brief: 单元测试-dao层-按用户Id更新数据
// @author: Kewin Li
// @receiver g
func (g *GormUserDaoSuite) TestUpdateById() {
	t := g.T()
	args := []driver.Value{}
	for i := 0; i < 7; i++ {
		args = append(args, sqlmock.AnyArg())
	}

	testCases := []struct {
		name string

		mock func(t *testing.T, u User) *sql.DB

		user User

		wantErr error
	}{
		// 查询成功
		{
			name: "Query successfully",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)

				mock.ExpectExec("UPDATE .*").
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)

				return db

			},
			user: User{
				Id:       1,
				Nickname: "test name",
				Phone:    sql.NullString{String: "111", Valid: true},
			},
		},
		// 没有记录
		{
			name: "Record not find",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("UPDATE .*").
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(gorm.ErrRecordNotFound)

				return db

			},
			user: User{
				Id:       2,
				Nickname: "test name",
				Phone:    sql.NullString{String: "222", Valid: true},
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		// 其他错误
		{
			name: "Other errors",
			mock: func(t *testing.T, u User) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("UPDATE .*").
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 0)).
					WillReturnError(errors.New("数据库错误"))

				return db

			},
			user: User{
				Id:       3,
				Nickname: "test name",
				Phone:    sql.NullString{String: "333", Valid: true},
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t, tc.user)
			defer sqlDB.Close()

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   true,
			})

			assert.NoError(t, err)
			d := NewGormUserDao(db)
			err = d.UpdateById(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)

		})

	}
}

func TestGormUserDao(t *testing.T) {
	suite.Run(t, &GormUserDaoSuite{})
}
