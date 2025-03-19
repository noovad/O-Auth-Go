package controller

import (
	"encoding/json"
	"errors"
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

func HandleGoogleLogin(ctx *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(config.OauthStateString)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleLogOut(c *gin.Context) {
	helper.DeleteTokens(c)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (controller *UsersAuthController) HandleGoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	content, err := controller.authService.GetUserInfo(state, code)
	if err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}
	var loginAccount map[string]interface{}
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	email, ok := loginAccount["email"].(string)
	if !ok {
		helper.InternalServerErrorResponse(ctx, helper.ErrFailedToGetEmail)
		return
	}

	var userId int
	userResponse, err := controller.usersService.FindByEmail(email)
	if errors.Is(err, helper.ErrUserNotFound) {
		userId, err = controller.usersService.CreateAndReturnID(data.CreateUsersRequest{
			Username: helper.GetUsernameFromEmail(email),
			Email:    email,
		})
		if err != nil {
			if errors.Is(err, helper.ErrFailedValidation) {
				helper.BadRequestResponse(ctx, err)
				return
			}

			helper.InternalServerErrorResponse(ctx, err)
			return
		}
	} else if err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	} else {
		userId = userResponse.Id
	}

	helper.CreateAccessToken(ctx, fmt.Sprintf("%d", userId))
	helper.CreateRefreshToken(ctx, fmt.Sprintf("%d", userId))

	ctx.Redirect(http.StatusFound, "/home")
}
