package router

import (
	"learn_o_auth-project/api"
	"learn_o_auth-project/api/controller"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/presentation"

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
