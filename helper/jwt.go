package helper

import (
	"go_auth-project/helper/responsejson"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(ctx *gin.Context, id string, secret string, duration time.Duration, cookieName string) error {
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

func CreateSignedToken(ctx *gin.Context, email string) error {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	return createToken(ctx, email, secret, time.Minute*1, "Signed-token")
}

func DeleteTokens(ctx *gin.Context) {
	ctx.SetCookie("Refresh-token", "", -1, "/", "", false, true)
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
}

func VerifySignedToken(ctx *gin.Context) (string, error) {
	signedToken, err := ctx.Cookie("Signed-token")
	if err != nil {
		return "", err
	}

	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims["id"].(string), nil
	} else {
		return "", ErrInvalidCredentials
	}
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
				responsejson.Unauthorized(ctx)
				ctx.Abort()
				return
			}
			id, ok := claims["id"]
			if !ok {
				responsejson.Unauthorized(ctx)
				ctx.Abort()
				return
			}
			CreateAccessToken(ctx, id.(string))
		} else {
			responsejson.Unauthorized(ctx)
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
		responsejson.Forbidden(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	ctx.Next()
}
