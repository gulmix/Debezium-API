package service

import (
	"context"
	"debezium_server/internal/models"
)

type UserRepository interface {
	Select(ctx context.Context, offset, limit int) ([]models.User, error)
}

type UserService struct {
	Repository UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repository: repo,
	}
}

func (s *UserService) GetUsers(ctx context.Context, offset, limit int) ([]models.User, error) {
	return s.Repository.Select(ctx, offset, limit)
}
