package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/database"
	"gophermart/pkg/logger"
	"net/http"
)

func GetOrdersHandler(db *sqlx.DB) http.HandlerFunc {
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

		orders, err := database.GetOrders(db, userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		logger.Sugar.Infow("Response orders to %s", userID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
