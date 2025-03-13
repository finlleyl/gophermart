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

func CreateUser(db *sqlx.DB, user models.User, loyalty models.LoyaltyAccount) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(
		`INSERT INTO users (id, login, password_hash, created_at)
		 VALUES (:id, :login, :password_hash, :created_at)`,
		user,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(
		`INSERT INTO loyalty_accounts (user_id, current_balance, withdrawn_balance)
		 VALUES (:user_id, :current_balance, :withdrawn_balance)`,
		loyalty,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(
		`INSERT INTO orders (id, user_id, number, status, accrual, uploaded_at)
         VALUES (:id, :user_id, :number, :status, :accrual, :uploaded_at)`,
		order,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(
		`UPDATE loyalty_accounts
         SET current_balance = current_balance + :accrual
         WHERE user_id = :user_id`,
		map[string]interface{}{
			"accrual": order.Accrual,
			"user_id": order.UserID,
		},
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetOrders(db *sqlx.DB, userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	query := `
		SELECT number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at ASC
	`
	if err := db.Select(&orders, query, userID); err != nil {
		return nil, err
	}
	return orders, nil
}

func GetBalance(db *sqlx.DB, userID uuid.UUID) (models.LoyaltyAccount, error) {
	var loyalty models.LoyaltyAccount

	query := `SELECT current_balance, withdrawn_balance FROM loyalty_accounts WHERE user_id = $1`
	if err := db.Get(&loyalty, query, userID); err != nil {
		return loyalty, err
	}

	return loyalty, nil
}

func Withdraw(db *sqlx.DB, orderNumber string, sum float64, userID uuid.UUID) error {
	var order models.Order
	query := `SELECT * FROM orders WHERE number = $1 AND user_id = $2`
	if err := db.Get(&order, query, orderNumber, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mErrors.ErrOrderNotFound
		}
		return err
	}

	var loyalty models.LoyaltyAccount
	if err := db.Get(&loyalty, `SELECT * FROM loyalty_accounts WHERE user_id = $1`, userID); err != nil {
		return err
	}

	if loyalty.CurrentBalance < sum {
		return mErrors.ErrNotEnoughMoney
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(
		`UPDATE loyalty_accounts
		 SET current_balance = current_balance - :sum, withdrawn_balance = withdrawn_balance + :sum
		 WHERE user_id = :user_id`,
		map[string]interface{}{
			"sum":     sum,
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(
		`INSERT INTO withdrawals (user_id, order_number, sum)
		 VALUES (:user_id, :order_number, :sum)`,
		map[string]interface{}{
			"user_id":      userID,
			"order_number": orderNumber,
			"sum":          sum,
		},
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
