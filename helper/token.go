package helper

import (
	"os"
	"time"

	"github.com/noovad/go-auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(claims jwt.MapClaims, secret string, duration time.Duration) string {
	now := time.Now()
	claims["exp"] = now.Add(duration).Unix()
	claims["iat"] = now.Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	return token
}

func CreateAccessToken(user dto.UserResponse) string {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id":          user.Id,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"avatar_type": user.AvatarType,
	}
	return generateToken(claims, secret, time.Minute*5)
}

func CreateRefreshToken(user dto.UserResponse) string {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id":          user.Id,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"avatar_type": user.AvatarType,
	}
	return generateToken(claims, secret, time.Hour*24*30)
}

func CreateSignedToken(email string) string {
	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id": email,
	}
	return generateToken(claims, secret, time.Minute*1)
}

func VerifySignedToken(ctx *gin.Context, email string) error {
	signedToken, _ := ctx.Cookie("Signed-token")

	secret := os.Getenv("GENERATE_SIGNING_TOKEN_SECRET")
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	SetCookie(ctx.Writer, "Signed-token", "", 0)

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
