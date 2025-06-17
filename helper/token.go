package helper

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(id string, secret string, duration time.Duration) string {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	return token
}

func CreateAccessToken(id string) string {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	return generateToken(id, secret, time.Minute*30)
}

func CreateRefreshToken(id string) string {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	return generateToken(id, secret, time.Hour*24*7)
}

func CreateSignedToken(email string) string {
	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	return generateToken(email, secret, time.Minute*1)
}

func VerifySignedToken(ctx *gin.Context, email string) error {
	signedToken, _ := ctx.Cookie("Signed-token")

	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return ErrInvalidCredentials
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if claims["id"] != email {
			return ErrInvalidCredentials
		}
		return nil
	} else {
		return ErrInvalidCredentials
	}
}
