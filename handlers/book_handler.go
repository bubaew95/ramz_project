package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ramz_project/models"
	"ramz_project/services"

	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{bookService: service}
}

// обработчик создания книги
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка парсинга JSON"))
		return
	}

	// Валидация данных
	if newBook.Name == "" || newBook.Year == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Название и год книги обязательны"))
		return
	}

	// Получаем user_id из куки
	cookie, err := r.Cookie("auth")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Авторизация требуется"))
		return
	}

	userId, err := strconv.Atoi(cookie.Value)
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неверная авторизация"))
		return
	}

	// Вызов сервиса
	bookId, err := h.bookService.CreateBook(r.Context(), &newBook, userId)
	if err != nil {
		log.Printf("Ошибка записи книги: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка записи в БД"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"book_id": bookId})
}

// обработчик получения всех книг
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books, err := h.bookService.GetAllBooks(r.Context())
	if err != nil {
		log.Printf("Ошибка получения книг: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка получения данных"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

// бработчик получения книги по ID
func (h *BookHandler) GetBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID"))
		return
	}

	book, err := h.bookService.GetBook(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Книга не найдена"))
		} else {
			log.Printf("Ошибка получения книги: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка получения данных"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// обработчик обновления книги
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID"))
		return
	}

	var updateBook models.Book
	err = json.NewDecoder(r.Body).Decode(&updateBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка парсинга JSON"))
		return
	}

	// Валидация данных
	if updateBook.Name == "" || updateBook.Year == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Название и год книги обязательны"))
		return
	}

	// Вызов сервиса
	err = h.bookService.UpdateBook(r.Context(), id, &updateBook)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Книга не найдена"))
		} else {
			log.Printf("Ошибка обновления книги: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка обновления"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateBook)
}

// обработчик удаления книги
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректный ID"))
		return
	}

	err = h.bookService.DeleteBook(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Книга не найдена"))
		} else {
			log.Printf("Ошибка удаления книги: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ошибка удаления"))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
