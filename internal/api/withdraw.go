package api

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/database"
	"gophermart/pkg/logger"
	"gophermart/pkg/mErrors"
	"net/http"
)

type WithdrawReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func WithdrawHandler(db *sqlx.DB) http.HandlerFunc {
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

		var req WithdrawReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err := database.Withdraw(db, req.Order, req.Sum, userID)
		if err != nil {
			if errors.Is(err, mErrors.ErrOrderNotFound) {
				http.Error(w, mErrors.ErrOrderNotFound.Error(), http.StatusUnprocessableEntity)
				return
			} else if errors.Is(err, mErrors.ErrNotEnoughMoney) {
				http.Error(w, mErrors.ErrNotEnoughMoney.Error(), http.StatusPaymentRequired)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Sugar.Infow("Withdraw from user %s for order %s with sum %f", userID, req.Order, req.Sum)
		w.WriteHeader(http.StatusOK)
	}
}
