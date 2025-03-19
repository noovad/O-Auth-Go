package router

import (
	"learn_o_auth-project/api"
	"learn_o_auth-project/api/controller"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/presentation"

	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	midleware := helper.RequireAccessToken

	r.GET("/", presentation.LoginPage)
	r.GET("/home", midleware, presentation.HomePage)
	r.GET("/login", controller.HandleGoogleLogin)
	r.GET("/logout", controller.HandleLogOut)
	r.GET("/callback", api.InitializeAuthController().HandleGoogleCallback)
}
