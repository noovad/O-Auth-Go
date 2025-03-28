package controller

import (
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/config"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/helper/responsejson"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersAuthController struct {
	usersService service.UsersService
	authService  service.AuthService
}

func NewUsersAuthController(userService service.UsersService, authService service.AuthService) *UsersAuthController {
	return &UsersAuthController{
		usersService: userService,
		authService:  authService,
	}
}

func HandleGoogleLogin(ctx *gin.Context) {
	state := helper.GenerateState()
	ctx.SetCookie("oauthstate", state, 60, "/", "", false, true)

	url := config.GoogleOauthConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleLogOut(c *gin.Context) {
	helper.DeleteTokens(c)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (controller *UsersAuthController) HandleGoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	user, err := controller.authService.AuthenticateWithGoogle(ctx, state, code)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	err = controller.authService.CreateTokens(ctx, user.Id)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/home")
}
