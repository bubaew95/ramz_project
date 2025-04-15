package repositories

import (
	"context"
	"ramz_project/models"
)

type BookRepository interface {
	Create(ctx context.Context, book *models.Book, userId int) (int, error)
	Update(ctx context.Context, id int, book *models.Book) error
	Get(ctx context.Context, id int) (*models.Book, error)
	GetAll(ctx context.Context) ([]*models.Book, error)
	Delete(ctx context.Context, id int) error
}
