package api

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRegisterHandler_ExistingUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "created_at"}).
		AddRow(uuid.New(), "existinguser", "hashedpassword", time.Now())

	mock.ExpectQuery("SELECT \\* FROM users WHERE login = \\$1").
		WithArgs("existinguser").
		WillReturnRows(rows)

	payload := `{"login": "existinguser", "password": "validpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(payload))
	w := httptest.NewRecorder()

	cfg := &config.Config{SecretKey: "testsecret"}
	handler := RegisterHandler(sqlxDB, cfg)
	handler(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "user already exists\n", w.Body.String())
}
