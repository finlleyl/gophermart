package user

import (
	"net/http"
	"time"
)

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
