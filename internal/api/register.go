package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"gophermart/pkg/hash"
	"gophermart/pkg/jwt"
	"gophermart/pkg/models"
	"net/http"
	"time"
)

type RequestRegistration struct {
	Login    string
	Password string
}

func RegisterHandler(db *sqlx.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestRegistration
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(req.Login) < 3 || len(req.Password) < 6 {
			http.Error(w, "invalid login or password", http.StatusBadRequest)
			return
		}

		if _, err := database.FindByLogin(db, req.Login); err == nil {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		hashedPassword, err := hash.Hash(req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			ID:           uuid.New(),
			Login:        req.Login,
			PasswordHash: hashedPassword,
			CreatedAt:    time.Now(),
		}

		if err := database.CreateUser(db, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateJWT(cfg, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(time.Hour * 3),
		})

		w.WriteHeader(http.StatusCreated)
	}
}
