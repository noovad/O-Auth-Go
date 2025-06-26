package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/noovad/go-auth/config"
	"github.com/noovad/go-auth/dto"
	"github.com/noovad/go-auth/helper/responsejson"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("access_token")
	user, valid := ValidateToken(accessToken, os.Getenv("GENERATE_TOKEN_SECRET"))
	if valid && ensureUserExists(ctx, user.Id.String()) {
		ctx.Set("userId", user.Id)
		ctx.Next()
		return
	}

	refreshToken, _ := ctx.Cookie("refresh_token")
	user, refreshValid := ValidateToken(refreshToken, os.Getenv("GENERATE_REFRESH_TOKEN_SECRET"))
	if refreshValid && ensureUserExists(ctx, user.Id.String()) {
		newAccessToken := CreateAccessToken(user)
		SetCookie(ctx.Writer, "access_token", newAccessToken, 60*60*24)
		ctx.Set("userId", user.Id)
		ctx.Next()
		return
	}

	responsejson.Unauthorized(ctx)
	ctx.Abort()
}

func GuestMiddleware(ctx *gin.Context) {
	accessToken, _ := ctx.Cookie("access_token")
	user, valid := ValidateToken(accessToken, os.Getenv("GENERATE_TOKEN_SECRET"))
	if valid && ensureUserExists(ctx, user.Id.String()) {
		responsejson.Forbidden(ctx, "You are already logged in")
		ctx.Abort()
		return
	}

	refreshToken, _ := ctx.Cookie("refresh_token")
	user, refreshValid := ValidateToken(refreshToken, os.Getenv("GENERATE_REFRESH_TOKEN_SECRET"))
	if refreshValid && ensureUserExists(ctx, user.Id.String()) {
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
		log.Printf("parseToken: Invalid token: %v", err)
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("parseToken: Invalid claims")
		return nil, nil, fmt.Errorf("invalid claims")
	}

	return token, claims, nil
}

func ValidateToken(tokenStr string, secret string) (dto.UserResponse, bool) {
	_, claims, err := parseToken(tokenStr, secret)
	if err != nil {
		log.Printf("ValidateToken: Token validation failed: %v", err)
		return dto.UserResponse{}, false
	}

	id, ok := claims["id"].(string)
	if !ok {
		log.Println("ValidateToken: ID claim missing or invalid")
		return dto.UserResponse{}, false
	}
	UserResponse := dto.UserResponse{
		Id:         uuid.MustParse(id),
		Name:       claims["name"].(string),
		Username:   claims["username"].(string),
		Email:      claims["email"].(string),
		AvatarType: claims["avatar_type"].(string),
	}
	return UserResponse, true
}

func ensureUserExists(ctx *gin.Context, userId string) bool {
	exists, err := UserExistsInDatabase(userId)
	if err != nil {
		log.Printf("ensureUserExists: Error checking user existence: %v", err)
		responsejson.InternalServerError(ctx, err)
		ctx.Abort()
		return false
	}
	if !exists {
		SetCookie(ctx.Writer, "access_token", "", -1)
		SetCookie(ctx.Writer, "refresh_token", "", -1)
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
		log.Printf("UserExistsInDatabase: Database error: %v", err)
		return false, err
	}

	return exists, nil
}
