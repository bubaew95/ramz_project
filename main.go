package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func init() {
	db = connectDB()
	migration()
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("auth")
		if err != nil {
			fmt.Println("test")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Вы не авторизованы"))
			return
		}

		h.ServeHTTP(w, r)
	})
}

// new commit 1

func main() {
	defer db.Close()
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/registration", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("route", "registration")
			//Регистрация добавление пользователя в БД

		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			type LoginRequest struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}

			var request LoginRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				log.Println("auth error", err)
			}

			//Проверка в БД на совпадения логина и пароля
			if request.Login == "admin" && request.Password == "123123" {
				cookie := &http.Cookie{
					Name:  "auth",
					Value: "100",
					Path:  "/",
				}
				http.SetCookie(w, cookie)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Авторизован"))
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Ошибка авторизации"))
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Get("/book/{id}", getBookById)
		r.Get("/books", getAllBooks)
		r.Post("/books", createBook)
		r.Put("/book/{id}", updateBook)
		r.Delete("/book/{id}", deleteBook)
	})

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server run error %v", err)
	}
}

type Book struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Year    int    `json:"year"`
	Image   string `json:"image"`
	Visible int    `json:"visible"`
}

func getBookById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if strings.TrimSpace(id) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is required"))
		return
	}

	// Разница между query и queryRow

	row := db.QueryRow(`SELECT id, name, year, image, visible FROM books WHERE id = $1`, id)
	if err := row.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка получения данных"))
		return
	}

	var book Book
	err := row.Scan(&book.Id, &book.Name, &book.Year, &book.Image, &book.Visible)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("В БД нет данных по id: " + id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte("Ошибка кодирования Json"))
		return
	}
}

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	rows, err := db.QueryContext(r.Context(), `SELECT id, name, year, image, visible FROM books`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка получения данных"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.Id, &book.Name, &book.Year, &book.Image, &book.Visible)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка сканирования данных"))
			return
		}
		books = append(books, book)
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка кодирования Json"))
	}
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newBook Book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка парсинга JSON"))
		return
	}

	if newBook.Name == "" || newBook.Year == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Названия и год книги обязательны"))
		return
	}

	err = insertBook(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка записи в БД"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID"))
		return
	}

	var updateBook Book
	err = json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка парсинга JSON"))
		return
	}

	if updateBook.Name == "" || updateBook.Year == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Названия и год книги обязательны"))
		return
	}

	err = updateBookInDB(id, &updateBook)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Book not found"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ошибка обновления"))
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateBook)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный id"))
		return
	}
	/// Почему?
	err = deleteFromDB(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Книга не найдена"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка удаления"))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Успешно удалена книга"))

}

func connectDB() *sql.DB {
	db, err := sql.Open("pgx", "host=127.0.0.1 user=admin password=admin dbname=ramz sslmode=disable")
	if err != nil {
		log.Fatalf("connect db error %v", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	// Ограничение на максимальное число открытых соединений
	db.SetMaxOpenConns(10)
	// Ограничение на количество соединений в простое
	db.SetMaxIdleConns(5)

	return db
}

func deleteFromDB(id int) error {
	query := `
        DELETE FROM books
        WHERE id = $1
    `
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func insertBook(book *Book) error {

	query := `
	    INSERT INTO books (name, year, image, visible)
    	VALUES ($1, $2, $3, $4)
    	RETURNING id
`

	err := db.QueryRow(query,
		book.Name,
		book.Year,
		book.Image,
		book.Visible,
	).Scan(&book.Id)

	if err != nil {
		return fmt.Errorf("Ошибка вставки книги: %w", err)
	}

	return nil

}

func updateBookInDB(id int, book *Book) error {

	query := `
	    UPDATE books
    	SET name = $1, year = $2, image = $3, visible = $4
    	WHERE id = $5
`

	res, err := db.Exec(query,
		book.Name,
		book.Year,
		book.Image,
		book.Visible,
		id+1,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func migration() {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
		    id serial PRIMARY KEY,
		    name VARCHAR(255) NOT NULL,
		    year INTEGER DEFAULT 0,
		    image VARCHAR(255) DEFAULT NULL,
		    visible INTEGER DEFAULT 1
		);
	`)
	if err != nil {
		log.Fatalf("migration err %v", err)
	}
}
