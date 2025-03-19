package controller

import (
	"encoding/json"
	"fmt"
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/config"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"
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

func HandleGoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(config.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleLogOut(c *gin.Context) {
	helper.DeleteTokens(c)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (controller *UsersAuthController) HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	content, err := controller.authService.GetUserInfo(state, code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	var loginAccount map[string]interface{}
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	email := loginAccount["email"].(string)

	userId := controller.usersService.FindByEmail(email).Id

	if userId == 0 {
		controller.usersService.Create(data.CreateUsersRequest{
			Username: helper.GetUsernameFromEmail(email),
			Email:    email,
		})
		userId = controller.usersService.FindByEmail(email).Id
	}

	helper.CreateAccessToken(c, fmt.Sprintf("%d", userId))
	helper.CreateRefreshToken(c, fmt.Sprintf("%d", userId))

	c.Redirect(http.StatusTemporaryRedirect, "/home")
}
