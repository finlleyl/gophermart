package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/pkg/config"
)

func TestLoginHandler_InvalidUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Эмулируем ситуацию, когда пользователь не найден
	mock.ExpectQuery("SELECT \\* FROM users WHERE login = \\$1").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	payload := `{"login": "nonexistent", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(payload))
	w := httptest.NewRecorder()

	cfg := &config.Config{SecretKey: "testsecret"}
	handler := LoginHandler(sqlxDB, cfg)
	handler(w, req)

	// Ожидаем, что статус ответа будет 401 (Unauthorized)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
