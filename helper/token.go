package helper

import (
	"strconv"
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
	secret := MustGetenv("GENERATE_ACCESS_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id":          user.Id,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"avatar_type": user.AvatarType,
	}
	ageDuration := MustParseDurationEnv("ACCESS_TOKEN_AGE")
	return generateToken(claims, secret, time.Second*time.Duration(ageDuration))
}

func CreateRefreshToken(user dto.UserResponse) string {
	secret := MustGetenv("GENERATE_REFRESH_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id":          user.Id,
		"name":        user.Name,
		"username":    user.Username,
		"email":       user.Email,
		"avatar_type": user.AvatarType,
	}
	ageDuration := MustParseDurationEnv("REFRESH_TOKEN_AGE")
	return generateToken(claims, secret, time.Second*time.Duration(ageDuration))
}

func CreateSignedToken(email string) string {
	secret := MustGetenv("GENERATE_SIGNING_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"id": email,
	}
	return generateToken(claims, secret, time.Minute*1)
}

func VerifySignedToken(ctx *gin.Context, email string) error {
	signedToken, _ := ctx.Cookie("Signed-token")

	secret := MustGetenv("GENERATE_SIGNING_TOKEN_SECRET")
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	SetCookie(ctx.Writer, "Signed-token", "", -1)

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

func MustParseDurationEnv(key string) int {
	value := MustGetenv(key)
	seconds, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid " + key + " format: " + err.Error())
	}
	return seconds
}
