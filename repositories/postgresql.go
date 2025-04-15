package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"ramz_project/models"
)

// реализует репозиторий для книг
type PostgreSQLBookRepository struct {
	db *sql.DB
}

// создает новый репозиторий книг
func NewPostgreSQLBookRepository(db *sql.DB) *PostgreSQLBookRepository {
	return &PostgreSQLBookRepository{db: db}
}

// добавляет новую книгу
func (r *PostgreSQLBookRepository) Create(ctx context.Context, book *models.Book, userId int) (int, error) {
	query := `
        INSERT INTO books (name, year, image, visible, user_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var newBookId int
	err := r.db.QueryRowContext(ctx, query,
		book.Name,
		book.Year,
		book.Image,
		book.Visible,
		userId,
	).Scan(&newBookId)

	if err != nil {
		return 0, fmt.Errorf("ошибка вставки книги: %w", err)
	}

	return newBookId, nil
}

// обновляет книгу
func (r *PostgreSQLBookRepository) Update(ctx context.Context, id int, book *models.Book) error {
	query := `
        UPDATE books
        SET name = $1, year = $2, image = $3, visible = $4
        WHERE id = $5
    `

	_, err := r.db.ExecContext(ctx, query,
		book.Name,
		book.Year,
		book.Image,
		book.Visible,
		id,
	)

	return err
}

// удаляет книгу
func (r *PostgreSQLBookRepository) Delete(ctx context.Context, id int) error {
	query := `
        DELETE FROM books
        WHERE id = $1
    `

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// получает книгу по ID
func (r *PostgreSQLBookRepository) Get(ctx context.Context, id int) (*models.Book, error) {
	query := `
        SELECT id, name, year, image, visible, user_id
        FROM books
        WHERE id = $1
    `

	var book models.Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Name,
		&book.Year,
		&book.Image,
		&book.Visible,
		&book.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("книга не найдена: %w", err)
		}
		return nil, fmt.Errorf("ошибка получения книги: %w", err)
	}

	return &book, nil
}

// получает все книги
func (r *PostgreSQLBookRepository) GetAll(ctx context.Context) ([]*models.Book, error) {
	query := `
        SELECT id, name, year, image, visible, user_id
        FROM books
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка книг: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Name,
			&book.Year,
			&book.Image,
			&book.Visible,
			&book.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования данных: %w", err)
		}
		books = append(books, &book)
	}

	return books, nil
}

// реализует репозиторий для пользователей
type PostgreSQLUserRepository struct {
	db *sql.DB
}

// создает новый репозиторий пользователей
func NewPostgreSQLUserRepository(db *sql.DB) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{db: db}
}

// добавляет нового пользователя
func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id
    `

	var newUserId int
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&newUserId)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки пользователя: %w", err)
	}

	return newUserId, nil
}

// получает пользователя по username
func (r *PostgreSQLUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
        SELECT id, username, password
        FROM users
        WHERE username = $1
    `

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return &user, nil
}
