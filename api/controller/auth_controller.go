package controller

import (
	"errors"
	"go_auth-project/api/service"
	"go_auth-project/config"
	"go_auth-project/dto"
	"go_auth-project/helper"
	"go_auth-project/helper/responsejson"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	usersService service.UserService
	authService  service.AuthService
}

func NewAuthController(userService service.UserService, authService service.AuthService) *AuthController {
	return &AuthController{
		usersService: userService,
		authService:  authService,
	}
}

func HandleGoogleAuth(ctx *gin.Context) {
	state := helper.GenerateState()
	helper.SetCookie(ctx, "Oauth-State", state, 60)

	url := config.GoogleOauthConfig.AuthCodeURL(state)
	responsejson.Success(ctx, "Redirecting to Google OAuth", gin.H{
		"url": url,
	})
}

func (c *AuthController) HandleSignUp(ctx *gin.Context) {
	email := ctx.Query("email")
	err := helper.VerifySignedToken(ctx, email)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	var user dto.CreateUsersRequest
	user.Email = email

	if err := ctx.ShouldBindJSON(&user); err != nil {
		responsejson.BadRequest(ctx, err)
		return
	}

	userId, err := c.usersService.CreateAndReturnID(user)
	if err != nil {
		if errors.Is(err, helper.ErrFailedValidation) {
			responsejson.BadRequest(ctx, err)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	err = c.authService.CreateTokens(ctx, userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "Successfully signed up", gin.H{
		"user":    user,
		"message": "You have been signed up successfully",
	})
}

func (c *AuthController) HandleLogin(ctx *gin.Context) {
	var user dto.LoginRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		responsejson.BadRequest(ctx, err)
		return
	}

	userId, err := c.authService.AuthenticateWithUsername(ctx, user)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	err = c.authService.CreateTokens(ctx, userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "Successfully logged in", gin.H{
		"user":    user,
		"message": "You have been logged in successfully",
	})
}

func HandleLogOut(c *gin.Context) {
	helper.SetCookie(c, "Authorization", "", -1)
	helper.SetCookie(c, "Refresh-token", "", -1)

	responsejson.Success(c, "Successfully logged out", gin.H{
		"message": "You have been logged out successfully",
	})
}

func (c *AuthController) HandleGoogleAuthCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	user, err := c.authService.AuthenticateWithGoogle(ctx, state, code)
	if err != nil {
		if errors.Is(err, helper.ErrOAuthStateNotFound) {
			ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=OAuth state not found")
			return
		}
		if errors.Is(err, helper.ErrInvalidOAuthState) {
			ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Invalid OAuth state")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.CreateSignedToken(ctx, user.Email)
			ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_BASE_URL")+"/sign-up?email="+user.Email)
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Failed to authenticate with Google")
		return
	}

	err = c.authService.CreateTokens(ctx, user.Id)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Failed to authenticate with Google")
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/")
}
