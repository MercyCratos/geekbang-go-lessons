package repository

import (
	"context"
	"geekbang-lessons/webook/internal/domain"
	"geekbang-lessons/webook/internal/repository/cache"
	"geekbang-lessons/webook/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
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
	du, err := repo.cache.Get(ctx, uid)
	// err 有两种可能：
	// 1. key not exist，说明 redis 是正常的
	// 2. 访问 redis 有问题；可能是网络有问题，也可能是 redis 本身就崩溃了
	if err == nil {
		return du, nil
	}

	u, err := repo.dao.SelectById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)

	// 异步回写缓存，其他语言可以用同步写，GO因为异步太容易写了，可以直接用异步写法，问题不大；
	go func() {
		err = repo.cache.Set(ctx, du)
		if err != nil {
			log.Println(err)
		}
	}()

	return du, nil
}

// FindByIdV1 对缓存处理比较较真的写法
func (repo *UserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// err 有两种可能：
	// 1. key not exist，说明 redis 是正常的
	// 2. 访问 redis 有问题；可能是网络有问题，也可能是 redis 本身就崩溃了
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		u, err := repo.dao.SelectById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)

		// 异步回写缓存，其他语言可以用同步写，GO因为异步太容易写了，可以直接用异步写法，问题不大；
		go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				log.Println(err)
			}
		}()

		return du, nil
	default:
		// redis 不正常
		// 高并发情况下，接近降级的保守写法
		return domain.User{}, err
	}
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
