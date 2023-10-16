package repository

import (
	"context"
	"geekbang-lessons/webook/internal/domain"
	"geekbang-lessons/webook/internal/repository/dao"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

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

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.SelectOneByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.dao.SelectById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toDataObject(user))
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *UserRepository) toDataObject(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
