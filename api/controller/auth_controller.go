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

func (controller *UsersAuthController) HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	content, err := controller.authService.GetUserInfo(state, code)
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	fmt.Println(string(content))
	var loginAccount map[string]interface{}
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	fmt.Println(loginAccount)

	email := loginAccount["email"].(string)

	userId := controller.usersService.FindByEmail(email).Id

	if userId == 0 {
		println(email)
		controller.usersService.Create(data.CreateUsersRequest{
			Username: helper.GetUsernameFromEmail(email),
			Email:    email,
		})
		userId = controller.usersService.FindByEmail(email).Id
	}

	token, err := helper.CreateAccessToken(fmt.Sprintf("%d", userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("Authorization", token, 60*30, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/home")
}
