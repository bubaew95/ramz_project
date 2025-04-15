package repositories

import (
	"context"
	"ramz_project/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}
