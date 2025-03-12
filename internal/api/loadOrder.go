package api

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/internal/service/luhn"
	"gophermart/pkg/database"
	"gophermart/pkg/logger"
	"gophermart/pkg/mErrors"
	"gophermart/pkg/models"
	"io"
	"net/http"
	"strings"
	"time"
)

const ctxUserIDKey = "userIDKey"

type ErrorResponse struct {
	Error string `json:"error"`
}

func LoadOrderHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Unable to read order number", http.StatusBadRequest)
			return
		}
		orderNumber := strings.TrimSpace(string(bodyBytes))
		if orderNumber == "" {
			http.Error(w, "Order number is empty", http.StatusBadRequest)
			return
		}

		for _, ch := range orderNumber {
			if ch < '0' || ch > '9' {
				http.Error(w, "Order number must contain only digits", http.StatusUnprocessableEntity)
				return
			}
		}

		if ok := luhn.CheckLuhn(orderNumber); !ok {
			http.Error(w, "Wrong order number", http.StatusUnprocessableEntity)
			return
		}

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

		err = database.CheckOrder(db, userID, orderNumber)
		if err != nil {
			switch err {
			case mErrors.ErrOrderAlreadyCreated:
				http.Error(w, err.Error(), http.StatusOK)
				return
			case mErrors.ErrOrderAlreadyUploaded:
				http.Error(w, err.Error(), http.StatusConflict)
				return
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		order := models.Order{
			ID:         uuid.NewString(),
			UserID:     userID,
			Number:     orderNumber,
			Status:     "PROCESSING",
			Accrual:    2.6,
			UploadedAt: time.Now(),
		}
		err = database.LoadOrder(db, order)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Sugar.Infof("Order %s accepted from %s", orderNumber, userID)
		w.WriteHeader(http.StatusAccepted)
	}
}
