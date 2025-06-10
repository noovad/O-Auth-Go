package router

import (
	"go_auth-project/api"
	"go_auth-project/api/controller"
	"go_auth-project/helper"
	"go_auth-project/presentation"

	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	authMidleware := helper.AuthMiddleware
	guestMiddleware := helper.GuestMiddleware

	r.GET("/", presentation.LoginPage)
	r.GET("/home", authMidleware, presentation.HomePage)
	r.GET("/login", guestMiddleware, controller.HandleGoogleLogin)
	r.GET("/logout", authMidleware, controller.HandleLogOut)
	r.GET("/callback", api.InitializeAuthController().HandleGoogleCallback)
}
