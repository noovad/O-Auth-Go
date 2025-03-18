package router

import (
	"learn_o_auth-project/api/controller"
	"learn_o_auth-project/presentation"

	"github.com/gin-gonic/gin"
)

func OAuthRoutes(r *gin.Engine) {
	r.GET("/", presentation.LoginPage)
	r.GET("/login", controller.HandleGoogleLogin)
	r.GET("/callback", controller.HandleGoogleCallback)
}
