package helper

import (
	"net/http"
	"os"
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   os.Getenv("BACKEND_DOMAIN"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSite(0),
		MaxAge:   maxAge,
	})
}
