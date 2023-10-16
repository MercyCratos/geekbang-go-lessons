package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrDuplicateEmail
		}
	}

	return err
}

func (dao *UserDao) SelectOneByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) UpdateById(ctx context.Context, entity User) error {
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"updated_at": time.Now().UnixMilli(),
			"nickname":   entity.Nickname,
			"birthday":   entity.Birthday,
			"about_me":   entity.AboutMe,
		}).Error
}

func (dao *UserDao) SelectById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 为什么用 Int? 避免时区问题
	// 这里永远使用 UTC 0 的毫秒数
	CreatedAt int64
	UpdatedAt int64
}
