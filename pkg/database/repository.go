package database

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

func FindByUserID(db *sqlx.DB, userID uuid.UUID) (models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateUser(db *sqlx.DB, user models.User) error {
	_, err := db.Exec("INSERT INTO users (id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Login, user.PasswordHash, user.CreatedAt)
	return err
}
