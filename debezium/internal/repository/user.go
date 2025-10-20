package repository

import (
	"context"
	"debezium_server/internal/models"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Select(ctx context.Context, offset, limit int) []models.User {
	return nil
}

func (r *UserRepository) SelectByID(ctx context.Context, id int64) models.User {
	return models.User{}
}

func (r *UserRepository) Insert(ctx context.Context, user models.User) error {
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user models.User) error {
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
