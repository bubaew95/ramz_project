package services

import (
	"context"
	"ramz_project/models"
	"ramz_project/repositories"
)

type BookService struct {
	bookRepo repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) *BookService {
	return &BookService{bookRepo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, book *models.Book, userId int) (int, error) {
	return s.bookRepo.Create(ctx, book, userId)
}

func (s *BookService) UpdateBook(ctx context.Context, id int, book *models.Book) error {
	return s.bookRepo.Update(ctx, id, book)
}

func (s *BookService) GetBook(ctx context.Context, id int) (*models.Book, error) {
	return s.bookRepo.Get(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*models.Book, error) {
	return s.bookRepo.GetAll(ctx)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, id)
}
