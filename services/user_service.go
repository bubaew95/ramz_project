package services

import (
	"context"

	"ramz_project/models"
	"ramz_project/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

// Create добавляет нового
func (s *UserService) Create(ctx context.Context, user *models.User) (int, error) {
	return s.userRepo.Create(ctx, user)
}

// GetByUsername получает пользователя
func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}
