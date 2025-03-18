package controller

import (
	"fmt"
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(config.OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	content, err := service.GetUserInfo(state, code)
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    string(content),
	})
}
