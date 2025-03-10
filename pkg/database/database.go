package database

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"os"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func loadDBConfig() *DBConfig {
	cfg := DBConfig{}

	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.User = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.Password = password
	}
	if dbname := os.Getenv("POSTGRES_DB"); dbname != "" {
		cfg.DBName = dbname
	}

	cfg.Port = 5432
	cfg.SSLMode = "disable"

	return &cfg
}

func ConnectDB() (*sqlx.DB, error) {
	cfg := loadDBConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
