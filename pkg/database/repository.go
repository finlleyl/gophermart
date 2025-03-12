package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gophermart/pkg/mErrors"
	"gophermart/pkg/models"
)

func FindByLogin(db *sqlx.DB, login string) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUser(db *sqlx.DB, user models.User) error {
	_, err := db.NamedExec("INSERT INTO users (id, login, password_hash, created_at) VALUES (:id, :login, :password_hash, :created_at)", user)
	return err
}

func CheckOrder(db *sqlx.DB, userID uuid.UUID, orderNumber string) error {
	var order models.Order

	query := `SELECT * FROM orders WHERE number = $1`
	err := db.Get(&order, query, orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	switch {
	case order.UserID == userID:
		return mErrors.ErrOrderAlreadyCreated
	default:
		return mErrors.ErrOrderAlreadyUploaded
	}
}

func LoadOrder(db *sqlx.DB, order models.Order) error {
	_, err := db.NamedExec("INSERT INTO orders (id, user_id, number, status, accrual, uploaded_at) VALUES (:id, :user_id, :number, :status, :accrual, :uploaded_at)", order)
	return err
}
