package service

import (
	"context"
	"errors"
	"go_learning/webook_project/webook/internal/domain"
	"go_learning/webook_project/webook/internal/repository"
	"go_learning/webook_project/webook/internal/repository/dao"

	"golang.org/x/crypto/bcrypt"
)

// 邮箱冲突的
var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

// 登录的
var ErrInvalidUserPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	//先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserPassword
	}
	return u, nil
}
