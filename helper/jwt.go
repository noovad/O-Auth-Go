package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(id string) (string, error) {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	var claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
