package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	user2 "gophermart/internal/service/user"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"gophermart/pkg/hash"
	"gophermart/pkg/jwt"
	"gophermart/pkg/logger"
	"gophermart/pkg/models"
	"net/http"
	"time"
)

type RequestRegistration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

		loyalty := models.LoyaltyAccount{
			UserID:           user.ID,
			CurrentBalance:   0,
			WithdrawnBalance: 0,
		}

		if err := database.CreateUser(db, user, loyalty); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := jwt.GenerateJWT(cfg, user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user2.SetJWT(w, token)

		logger.Sugar.Infow("user registered", "login", req.Login)

		w.WriteHeader(http.StatusCreated)
	}
}
