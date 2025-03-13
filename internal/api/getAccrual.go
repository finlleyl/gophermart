package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/database"
	"gophermart/pkg/mErrors"
	"net/http"
)

func GetAccrualHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderNumber := chi.URLParam(r, "number")
		order, err := database.GetAccrual(db, orderNumber)
		if err != nil {
			if errors.Is(err, mErrors.ErrOrderNotFound) {
				http.Error(w, "Order not found", http.StatusNoContent)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

}
