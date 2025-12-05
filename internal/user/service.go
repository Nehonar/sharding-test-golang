package user

import (
	"context"
	"go-sharding-basic/internal/models"
	"go-sharding-basic/internal/storage/router"
)

type UserService interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUser(ctx context.Context, username string) (*models.User, error)
}

type userService struct {
	storage router.UserStorage
}

func NewUserService(storage router.UserStorage) UserService {
	return &userService{
		storage: storage,
	}
}

func (s *userService) CreateUser(ctx context.Context, username string, password string) error {
	err := s.storage.SaveUser(ctx, username, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetUser(ctx context.Context, username string) (*models.User, error) {
	return s.storage.GetUser(ctx, username)
}
