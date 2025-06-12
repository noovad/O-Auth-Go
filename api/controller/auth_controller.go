package controller

import (
	"errors"
	"go_auth-project/api/service"
	"go_auth-project/config"
	"go_auth-project/data"
	"go_auth-project/helper"
	"go_auth-project/helper/responsejson"

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

func (controller *UsersAuthController) HandleSignUp(ctx *gin.Context) {
	email, err := helper.VerifySignedToken(ctx)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	var user data.CreateUsersRequest
	user.Email = email

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responsejson.BadRequest(ctx, err)
		return
	}

	userId, err := controller.usersService.CreateAndReturnID(user)
	if err != nil {
		if errors.Is(err, helper.ErrFailedValidation) {
			responsejson.BadRequest(ctx, err)
			return
		}

		responsejson.InternalServerError(ctx, err)
		return
	}

	err = controller.authService.CreateTokens(ctx, userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "Successfully signed up", gin.H{
		"user":    user,
		"message": "You have been signed up successfully",
	})
}

func (controller *UsersAuthController) HandleLogin(ctx *gin.Context) {
	var user data.LoginRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		responsejson.BadRequest(ctx, err)
		return
	}

	userId, err := controller.authService.AuthenticateWithPassword(ctx, user)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	err = controller.authService.CreateTokens(ctx, userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "Successfully logged in", gin.H{
		"user":    user,
		"message": "You have been logged in successfully",
	})
}

func HandleGoogleAuth(ctx *gin.Context) {
	action := ctx.Query("action")
	if action != "signup" {
		action = "login"
	}

	state := helper.GenerateState() + "|" + action

	ctx.SetCookie("oauthstate", state, 60, "/", "", false, true)

	url := config.GoogleOauthConfig.AuthCodeURL(state)
	responsejson.Success(ctx, "Redirecting to Google OAuth", gin.H{
		"url": url,
	})
}

func HandleLogOut(c *gin.Context) {
	helper.DeleteTokens(c)

	responsejson.Success(c, "Successfully logged out", gin.H{
		"message": "You have been logged out successfully",
	})
}

func (controller *UsersAuthController) HandleGoogleAuthCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	user, action, err := controller.authService.AuthenticateWithGoogle(ctx, state, code)

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) && action == "signup" {
			token := helper.CreateSignedToken(ctx, user.Email)
			responsejson.Success(ctx, "Redirecting to sign-up page", gin.H{
				"token": token,
			})
			return
		}

		if errors.Is(err, helper.ErrUserNotFound) {
			responsejson.Unauthorized(ctx)
			return
		}

		responsejson.InternalServerError(ctx, err)
		return
	}

	if action == "signup" {
		responsejson.Conflict(ctx, "User already exists, please log in instead")
		return
	}

	err = controller.authService.CreateTokens(ctx, user.Id)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "Successfully logged in", gin.H{
		"user":    user,
		"message": "You have been logged in successfully",
	})
}
