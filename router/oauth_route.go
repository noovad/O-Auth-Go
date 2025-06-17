package router

import (
	"go_auth-project/api"
	"go_auth-project/api/controller"
	"go_auth-project/helper"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	authMidleware := helper.AuthMiddleware
	guestMiddleware := helper.GuestMiddleware
	authController := api.AuthInjector()

	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Refresh-token", "Signed-token", "Oauth-State"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/sign-up", guestMiddleware, authController.HandleSignUp)
	r.POST("/login", guestMiddleware, authController.HandleLogin)
	r.POST("/logout", authMidleware, controller.HandleLogOut)
	r.GET("/auth", guestMiddleware, controller.HandleGoogleAuth)
	r.GET("/callback", authController.HandleGoogleAuthCallback)
}
