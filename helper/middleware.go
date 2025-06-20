package helper

import (
	"context"
	"fmt"
	"go_auth-project/config"
	"go_auth-project/helper/responsejson"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("access_token")
	userId, valid := ValidateToken(accessToken, os.Getenv("GENERATE_TOKEN_SECRET"))
	if valid && ensureUserExists(ctx, userId) {
		ctx.Set("userId", userId)
		ctx.Next()
		return
	}

	refreshToken, _ := ctx.Cookie("refresh_token")
	refreshUserId, refreshValid := ValidateToken(refreshToken, os.Getenv("GENERATE_REFRESH_TOKEN_SECRET"))
	if refreshValid && ensureUserExists(ctx, refreshUserId) {
		newAccessToken := CreateAccessToken(refreshUserId)
		SetCookie(ctx.Writer, "access_token", newAccessToken, 60*60*24)
		ctx.Set("userId", refreshUserId)
		ctx.Next()
		return
	}

	responsejson.Unauthorized(ctx)
	ctx.Abort()
}

func GuestMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("access_token")
	userId, valid := ValidateToken(accessToken, os.Getenv("GENERATE_TOKEN_SECRET"))
	if valid && ensureUserExists(ctx, userId) {
		responsejson.Forbidden(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	refreshToken, _ := ctx.Cookie("refresh_token")
	refreshUserId, refreshValid := ValidateToken(refreshToken, os.Getenv("GENERATE_REFRESH_TOKEN_SECRET"))
	if refreshValid && ensureUserExists(ctx, refreshUserId) {
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

	return token, claims, nil
}

func ValidateToken(tokenStr string, secret string) (string, bool) {
	_, claims, err := parseToken(tokenStr, secret)
	if err != nil {
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
