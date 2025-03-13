package jwt

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gophermart/pkg/config"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	cfg := &config.Config{SecretKey: "testsecret"}
	userID := uuid.New()

	token, err := GenerateJWT(cfg, userID)
	require.NoError(t, err, "Ошибка генерации токена")

	parsedUserID, err := GetUserID(cfg, token)
	require.NoError(t, err, "Ошибка парсинга токена")
	require.Equal(t, userID, parsedUserID, "UserID должен совпадать с исходным")
}
