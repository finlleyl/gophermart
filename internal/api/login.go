package api

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	user2 "gophermart/internal/service/user"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"gophermart/pkg/hash"
	"gophermart/pkg/jwt"
	"net/http"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginHandler(db *sqlx.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(req.Login) < 3 || len(req.Password) < 6 {
			http.Error(w, "invalid login or password", http.StatusBadRequest)
			return
		}

		user, err := database.FindByLogin(db, req.Login)
		if err != nil {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		}

		if err := hash.CheckPassword(user.PasswordHash, req.Password); err != nil {
			http.Error(w, "Invalid login or password", http.StatusUnauthorized)
			return
		}

		token, err := jwt.GenerateJWT(cfg, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user2.SetJWT(w, token)

		w.WriteHeader(http.StatusOK)
	}
}
