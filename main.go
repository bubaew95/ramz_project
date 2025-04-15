package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ramz_project/middleware"
	"time"

	"ramz_project/handlers"
	"ramz_project/repositories"
	"ramz_project/services"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func init() {
	db = connectDB()
	migration()
}

func connectDB() *sql.DB {
	db, err := sql.Open("pgx", "host=127.0.0.1 user=postgres password=your_password dbname=bookstore sslmode=disable")
	if err != nil {
		log.Fatalf("connect db error %v", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db
}

func migration() {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL
        );
        CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            year INTEGER DEFAULT 0,
            image VARCHAR(255),
            visible INTEGER DEFAULT 1,
            user_id INTEGER REFERENCES users(id) ON DELETE SET NULL
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	fmt.Println("Миграция успешно выполнена!")
}

func main() {
	defer db.Close()

	// репозиторие
	bookRepo := repositories.NewPostgreSQLBookRepository(db)
	userRepo := repositories.NewPostgreSQLUserRepository(db)

	// сервис
	bookService := services.NewBookService(bookRepo)
	userService := services.NewUserService(userRepo)

	// обработчики
	bookHandler := handlers.NewBookHandler(bookService)
	userHandler := handlers.NewUserHandler(userService)

	router := chi.NewRouter()

	router.Post("/registration", userHandler.RegisterUser)
	router.Post("/login", userHandler.LoginUser)

	protected := chi.NewRouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Post("/books", bookHandler.CreateBook)
	protected.Get("/books", bookHandler.GetAllBooks)
	protected.Get("/book/{id}", bookHandler.GetBookById)
	protected.Put("/book/{id}", bookHandler.UpdateBook)
	protected.Delete("/book/{id}", bookHandler.DeleteBook)

	router.Mount("/api", protected)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server run error %v", err)
	}
}
