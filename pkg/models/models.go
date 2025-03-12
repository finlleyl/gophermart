package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Login        string    `db:"login" json:"login"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Order struct {
	ID         string    `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	Number     string    `db:"number" json:"number"`
	Status     string    `db:"status" json:"status"`
	Accrual    float64   `db:"accrual" json:"accrual"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}
