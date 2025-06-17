package helper

import (
	"fmt"
	"go_auth-project/helper/responsejson"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetCookie(ctx *gin.Context, name, value string, duration time.Duration) {
	domain := "o-auth-go-production.up.railway.app"
	ctx.SetCookie(name, value, int(duration.Seconds()), "/", domain, true, false)
	// Explicitly add SameSite=None since Gin does not support it directly
	ctx.Writer.Header().Add("Set-Cookie",
		fmt.Sprintf("%s=%s; Path=/; Max-Age=%d; Domain=%s; Secure; HttpOnly; SameSite=None",
			name,
			value,
			int(duration.Seconds()),
			domain,
		),
	)
}

func generateToken(ctx *gin.Context, id string, secret string, duration time.Duration, cookieName string) error {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return err
	}

	SetCookie(ctx, cookieName, token, duration)
	return nil
}

func CreateAccessToken(ctx *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	return generateToken(ctx, id, secret, time.Minute*30, "Authorization")
}

func CreateRefreshToken(ctx *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	return generateToken(ctx, id, secret, time.Hour*24*7, "Refresh-token")
}

func CreateSignedToken(ctx *gin.Context, email string) error {
	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	return generateToken(ctx, email, secret, time.Minute*1, "Signed-token")
}

func VerifySignedToken(ctx *gin.Context, email string) error {
	signedToken, err := ctx.Cookie("Signed-token")
	if err != nil {
		return err
	}

	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if claims["id"] != email {
			return ErrInvalidCredentials
		}
		return nil
	} else {
		return ErrInvalidCredentials
	}
}
