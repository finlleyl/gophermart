package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecret"
	hashed, err := Hash(password)
	require.NoError(t, err, "Ошибка при хешировании пароля")

	// Проверка, что правильный пароль проходит проверку
	err = CheckPassword(hashed, password)
	require.NoError(t, err, "Пароль должен соответствовать хешу")

	// Проверка, что неверный пароль не проходит проверку
	err = CheckPassword(hashed, "wrongpassword")
	require.Error(t, err, "Неверный пароль не должен соответствовать хешу")
}
