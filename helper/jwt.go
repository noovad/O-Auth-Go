package helper

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(c *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	var claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	}

	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtClaim.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return err
	}

	c.SetCookie("Authorization", token, 60*30, "/", "", false, true)
	return nil
}

func CreateRefreshToken(c *gin.Context, id string) error {
	secret := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	var claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	jwtClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtClaim.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return err
	}

	c.SetCookie("Refresh-token", token, 60*60*24*7, "/", "", false, true)
	return nil
}

func DeleteTokens(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.SetCookie("Refresh-token", "", -1, "/", "", false, true)
}

func RequireAccessToken(c *gin.Context) {
	accessToken, _ := c.Cookie("Authorization")
	secret := os.Getenv("GENERATE_TOKEN_SECRET")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	refresToken, _ := c.Cookie("Refresh-token")
	secretRefresh := os.Getenv("GENERATE_REFRESH_TOKEN_SECRET")
	refreshToken, _ := jwt.Parse(refresToken, func(token *jwt.Token) (any, error) {
		return []byte(secretRefresh), nil
	})

	if err != nil || !token.Valid {
		if refresToken != "" && refreshToken.Valid {
			claims, ok := refreshToken.Claims.(jwt.MapClaims)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			id := claims["id"].(string)
			fmt.Println(id)
			fmt.Println(claims)
			CreateAccessToken(c, id)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
