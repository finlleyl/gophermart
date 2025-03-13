package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/database"
	"gophermart/pkg/logger"
	"net/http"
)

func GetBalanceHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDVal := r.Context().Value(ctxUserIDKey)
		if userIDVal == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID, ok := userIDVal.(uuid.UUID)
		if !ok {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		loyalty, err := database.GetBalance(db, userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(loyalty); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Sugar.Infow("Response balance to %s", userID)
	}
}
