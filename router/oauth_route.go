package router

import (
	"go_auth-project/api"
	"go_auth-project/api/controller"
	"go_auth-project/helper"

	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	authMidleware := helper.AuthMiddleware
	guestMiddleware := helper.GuestMiddleware
	authController := api.InitializeAuthController()

	r.POST("/sign-up", guestMiddleware, authController.HandleSignUp)
	r.POST("/login", guestMiddleware, authController.HandleLogin)
	r.GET("/logout", authMidleware, controller.HandleLogOut)
	r.GET("/auth", guestMiddleware, controller.HandleGoogleAuth)
	r.GET("/callback", authController.HandleGoogleAuthCallback)
}
