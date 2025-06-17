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
	ctx.SetCookie("Oauth-State", state, 60, "/", os.Getenv("BACKEND_DOMAIN"), true, true)

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

	responsejson.Success(ctx, "Successfully signed up", gin.H{
		"AccessToken":  helper.CreateAccessToken(userId.String()),
		"RefreshToken": helper.CreateRefreshToken(userId.String()),
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

	responsejson.Success(ctx, "Successfully logged in", gin.H{
		"AccessToken":  helper.CreateAccessToken(userId.String()),
		"RefreshToken": helper.CreateRefreshToken(userId.String()),
	})
}

func HandleLogOut(c *gin.Context) {

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
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=OAuth state not found")
			return
		}
		if errors.Is(err, helper.ErrInvalidOAuthState) {
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Invalid OAuth state")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.SetCookie("Signed-token", helper.CreateSignedToken(user.Email), 60*5, "/", os.Getenv("BACKEND_DOMAIN"), true, true)
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/sign-up?email="+user.Email)
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Failed to authenticate with Google")
		return
	}

	accessToken := helper.CreateAccessToken(user.Id.String())
	refreshToken := helper.CreateRefreshToken(user.Id.String())
	redirectURL := os.Getenv("FRONTEND_BASE_URL") + "/callback?accessToken=" + accessToken + "&refreshToken=" + refreshToken
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)

}
