package repository

import (
	"context"
	"geekbang-lessons/webook/internal/domain"
	"geekbang-lessons/webook/internal/repository/dao"
)

var ErrDuplicateEmail = dao.ErrDuplicateEmail

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
