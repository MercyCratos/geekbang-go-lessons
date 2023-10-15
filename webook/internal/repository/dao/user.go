package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
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

	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 为什么用 Int? 避免时区问题
	// 这里永远使用 UTC 0 的毫秒数
	CreatedAt int64
	UpdatedAt int64
}
