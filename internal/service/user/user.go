package user

import (
	"gophermart/pkg/config"
	"gophermart/pkg/jwt"
	"net/http"
	"time"
)

func CheckCookies(cfg *config.Config, r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	userID, err := jwt.GetUserID(cfg, tokenCookie.Value)
	if err != nil {
		return "", err
	}

	return userID.String(), nil
}

func SetJWT(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 3),
	})
}
