package middleware

import (
	"context"
	"gophermart/pkg/config"
	"gophermart/pkg/jwt"
	"net/http"
)

const ctxUserIDKey = "userIDKey"

func CheckCookies(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		userID, err := jwt.GetUserID(cfg, cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
