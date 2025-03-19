package helper

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(ctx *gin.Context, id, secret string, duration time.Duration, cookieName string) error {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return err
	}

	ctx.SetCookie(cookieName, token, int(duration.Seconds()), "/", "", false, true)
	return nil
}

func CreateAccessToken(ctx *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	return createToken(ctx, id, secret, time.Minute*30, "Authorization")
}

func CreateRefreshToken(ctx *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	return createToken(ctx, id, secret, time.Hour*24*7, "Refresh-token")
}

func DeleteTokens(ctx *gin.Context) {
	ctx.SetCookie("Refresh-token", "", -1, "/", "", false, true)
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
}

func AuthMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("Authorization")
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	refresToken, _ := ctx.Cookie("Refresh-token")
	secretRefresh := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	refreshToken, _ := jwt.Parse(refresToken, func(token *jwt.Token) (any, error) {
		return []byte(secretRefresh), nil
	})

	if err != nil || !token.Valid {
		if refresToken != "" && refreshToken.Valid {
			claims, ok := refreshToken.Claims.(jwt.MapClaims)
			if !ok {
				UnauthorizedResponse(ctx)
				ctx.Abort()
				return
			}
			id := claims["id"].(string)
			CreateAccessToken(ctx, id)
		} else {
			UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}
	}

	ctx.Next()
}

func GuestMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("Authorization")
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	refresToken, _ := ctx.Cookie("Refresh-token")
	secretRefresh := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	refreshToken, _ := jwt.Parse(refresToken, func(token *jwt.Token) (any, error) {
		return []byte(secretRefresh), nil
	})

	if (err == nil && accessToken != "" && token.Valid) || (err == nil && refresToken != "" && refreshToken.Valid) {
		ForbiddenResponse(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	ctx.Next()
}
