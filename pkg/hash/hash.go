package hash

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashed_password string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
}
