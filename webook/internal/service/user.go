package service

import (
	"context"
	"errors"
	"geekbang-lessons/webook/internal/domain"
	"geekbang-lessons/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(encrypted)

	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, ErrInvalidUserOrPassword
		}
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *UserService) UpdateUserInfo(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *UserService) GetProfile(ctx *gin.Context, userId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userId)
}
