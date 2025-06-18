package helper

import (
	"context"
	"fmt"
	"go_auth-project/config"
	"go_auth-project/helper/responsejson"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	accessToken := AccessTokenFromHeader(ctx)
	userId, valid := ValidateAccessToken(accessToken)
	if valid && ensureUserExists(ctx, userId) {
		ctx.Set("userId", userId)
		ctx.Next()
		return
	}

	responsejson.Unauthorized(ctx, "Invalid or expired access token")
	ctx.Abort()
}

func GuestMiddleware(ctx *gin.Context) {
	accessToken := AccessTokenFromHeader(ctx)

	userId, valid := ValidateAccessToken(accessToken)
	if valid && ensureUserExists(ctx, userId) {
		responsejson.Forbidden(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	ctx.Next()
}

func parseToken(tokenStr, secret string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
		return nil, nil, fmt.Errorf("token expired")
	}

	return token, claims, nil
}

func ValidateAccessToken(tokenStr string) (string, bool) {
	_, claims, err := parseToken(tokenStr, os.Getenv("GENERATE_TOKEN_SECRET"))
	if err != nil {
		fmt.Println("Invalid access token:", err)
		return "", false
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", false
	}
	return id, true
}

func ensureUserExists(ctx *gin.Context, userId string) bool {
	exists, err := UserExistsInDatabase(userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		ctx.Abort()
		return false
	}
	if !exists {
		responsejson.Unauthorized(ctx, "User not found")
		ctx.Abort()
		return false
	}
	return true
}

func UserExistsInDatabase(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := config.DatabaseConnection()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	err := db.WithContext(ctx).Raw(query, userId).Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func ValidateRefreshToken(tokenStr string) (string, bool) {
	_, claims, err := parseToken(tokenStr, os.Getenv("GENERATE_REFRESH_TOKEN_SECRET"))
	if err != nil {
		fmt.Println("Invalid refresh token:", err)
		return "", false
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", false
	}
	return id, true
}

func AccessTokenFromHeader(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
