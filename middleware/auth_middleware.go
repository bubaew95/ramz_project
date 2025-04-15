package middleware

import (
	"context"
	"net/http"
	"strconv"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Авторизация требуется", http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "Неверная авторизация", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
