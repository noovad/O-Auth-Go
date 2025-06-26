package router

import (
	"os"
	"time"

	"github.com/noovad/go-auth/api"
	"github.com/noovad/go-auth/api/controller"
	"github.com/noovad/go-auth/helper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	guestMiddleware := helper.GuestMiddleware
	authMiddleware := helper.AuthMiddleware
	authController := api.AuthInjector()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_BASE_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Refresh-token", "Signed-token", "Oauth-State"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	{
		auth := r.Group("/auth")
		auth.POST("/sign-up", guestMiddleware, authController.HandleSignUp)
		auth.POST("/login", guestMiddleware, authController.HandleLogin)
		auth.POST("/logout", authController.HandleLogout)
		auth.GET("/google", guestMiddleware, controller.HandleGoogleAuth)
		auth.GET("/callback", authController.HandleGoogleAuthCallback)
		auth.DELETE("/delete-account", authMiddleware, authController.HandleDeleteAccount)
		auth.PUT("/update-avatar", authMiddleware, authController.HandleUpdateAvatar)
	}
}
