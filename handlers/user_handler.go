package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ramz_project/models"
	"ramz_project/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

// обработчик регистрации
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ошибка парсинга JSON"))
		return
	}

	if reqData.Username == "" || reqData.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Логин и пароль обязательны"))
		return
	}

	user := &models.User{
		Username: reqData.Username,
		Password: reqData.Password,
	}

	// Проверка уникальности username
	existsUser, err := h.userService.GetByUsername(r.Context(), reqData.Username)
	if err == nil && existsUser != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Логин занят"))
		return
	}

	// Создаем пользователя
	newUserId, err := h.userService.Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ошибка сохранения"))
		return
	}

	// Установка куки
	cookie := &http.Cookie{
		Name:  "auth",
		Value: fmt.Sprintf("%d", newUserId),
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"user_id": newUserId})
}

// обработчик входа
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ошибка парсинга", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetByUsername(r.Context(), req.Username)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	if user.Password != req.Password {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:  "auth",
		Value: fmt.Sprintf("%d", user.ID),
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Авторизация успешна"))
}
