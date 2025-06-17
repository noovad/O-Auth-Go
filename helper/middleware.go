package helper

import (
	"go_auth-project/helper/responsejson"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
