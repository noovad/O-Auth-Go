package helper

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var response = Response{}

func createToken(ctx *gin.Context, id int, secret string, duration time.Duration, cookieName string) error {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		response.InternalServerError(ctx, err)
		return err
	}

	ctx.SetCookie(cookieName, token, int(duration.Seconds()), "/", "", false, true)
	return nil
}

func CreateAccessToken(ctx *gin.Context, id int) error {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	return createToken(ctx, id, secret, time.Minute*30, "Authorization")
}

func CreateRefreshToken(ctx *gin.Context, id int) error {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	return createToken(ctx, id, secret, time.Hour*24*7, "Refresh-token")
}

func DeleteTokens(ctx *gin.Context) {
	ctx.SetCookie("Refresh-token", "", -1, "/", "", false, true)
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
}

// Note: This middleware can be enhanced by:
// - Checking the token expiration time explicitly.
// - Validating the user ID from the token with the database.
// - Implementing role-based access control (RBAC) for different user permissions.
// - Storing the authenticated user in the request context using ctx.Set("user", user) for easier access in handlers.
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
				response.Unauthorized(ctx)
				ctx.Abort()
				return
			}
			idFloat, ok := claims["id"].(float64)
			if !ok {
				response.Unauthorized(ctx)
				ctx.Abort()
				return
			}
			id := int(idFloat)
			CreateAccessToken(ctx, id)
		} else {
			response.Unauthorized(ctx)
			ctx.Abort()
			return
		}
	}

	ctx.Next()
}

// Note: This middleware can be enhanced by:
// - Ensuring that both access and refresh tokens are properly invalidated when logging out.
// - Redirecting authenticated users away from guest-only pages instead of returning a Forbidden response.
// - Storing the guest status in the request context using ctx.Set("user", nil) for consistency.
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
		response.Forbidden(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	ctx.Next()
}
