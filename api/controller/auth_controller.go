package controller

import (
	"errors"
	"fmt"
	"go_auth-project/api/service"
	"go_auth-project/config"
	"go_auth-project/dto"
	"go_auth-project/helper"
	"go_auth-project/helper/responsejson"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthController struct {
	userService service.UserService
	authService service.AuthService
}

func NewAuthController(userService service.UserService, authService service.AuthService) *AuthController {
	return &AuthController{
		userService: userService,
		authService: authService,
	}
}

func HandleGoogleAuth(ctx *gin.Context) {
	state := helper.GenerateState()
	helper.SetCookie(ctx.Writer, "Oauth-State", state, 60)

	url := config.GoogleOauthConfig.AuthCodeURL(state)

	responsejson.Success(ctx, gin.H{
		"url": url,
	}, "Redirecting to Google OAuth")
}

func (c *AuthController) HandleSignUp(ctx *gin.Context) {
	email := ctx.Query("email")
	err := helper.VerifySignedToken(ctx, email)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx)
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to verify signed token")
		return
	}

	var userReq dto.CreateUsersRequest
	userReq.Email = email
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		responsejson.BadRequest(ctx, err, "Invalid sign-up request")
		return
	}

	user, err := c.userService.CreateUser(userReq)
	if err != nil {
		if errors.Is(err, helper.ErrFailedValidation) {
			responsejson.BadRequest(ctx, err, "Failed to create user due to validation error")
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to create user")
		return
	}

	refreshToken := helper.CreateRefreshToken(user)
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, 60*60*24*30)

	accessToken := helper.CreateAccessToken(user)
	helper.SetCookie(ctx.Writer, "access_token", accessToken, 60*5)
	responsejson.Success(ctx, gin.H{
		"user": user,
	}, "Successfully signed up")
}

func (c *AuthController) HandleLogin(ctx *gin.Context) {
	var user dto.LoginRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		responsejson.BadRequest(ctx, err, "Invalid login request")
		return
	}

	userResponse, err := c.authService.AuthenticateWithUsername(ctx, user)
	if err != nil {
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx, "Invalid username or password")
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to authenticate user")
		return
	}

	refreshToken := helper.CreateRefreshToken(userResponse)
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, 60*60*24*30)

	accessToken := helper.CreateAccessToken(userResponse)
	helper.SetCookie(ctx.Writer, "access_token", accessToken, 60*5)

	responsejson.Success(ctx, gin.H{
		"user": userResponse,
	}, "Successfully logged in")
}

func (c *AuthController) HandleLogout(ctx *gin.Context) {
	helper.SetCookie(ctx.Writer, "refresh_token", "", -1)
	helper.SetCookie(ctx.Writer, "access_token", "", -1)
	responsejson.Success(ctx, nil, "Successfully logged out")
}

func (c *AuthController) HandleGoogleAuthCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")

	user, err := c.authService.AuthenticateWithGoogle(ctx, state, code)
	if err != nil {
		if errors.Is(err, helper.ErrOAuthStateNotFound) {
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Invalid credentials, please login again")
			return
		}
		if errors.Is(err, helper.ErrInvalidOAuthState) {
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Invalid credentials, please login again")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.SetCookie("Signed-token", helper.CreateSignedToken(user.Email), 60*5, "/", os.Getenv("BACKEND_DOMAIN"), true, true)
			ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/sign-up?email="+user.Email)
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, os.Getenv("FRONTEND_BASE_URL")+"/login?error=Failed to authenticate with Google. Please try again later.")
		return
	}

	refreshToken := helper.CreateRefreshToken(user)
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, 60*60*24*30)
	accessToken := helper.CreateAccessToken(user)
	helper.SetCookie(ctx.Writer, "access_token", accessToken, 60*5)

	redirectURL := os.Getenv("FRONTEND_BASE_URL") + "/"
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (c *AuthController) HandleDeleteAccount(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	fmt.Println("HandleDeleteAccount: userId from context =", userId, "exists =", exists) // debug
	if !exists {
		responsejson.InternalServerError(ctx, nil, "User ID not found in context")
		return
	}
	uid, ok := userId.(uuid.UUID)
	if !ok {
		responsejson.InternalServerError(ctx, nil, "Invalid user ID type")
		return
	}
	err := c.userService.DeleteById(uid)
	fmt.Println("HandleDeleteAccount: error from DeleteById =", err) // debug
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responsejson.NotFound(ctx, "User not found")
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to delete user")
		return
	}

	helper.SetCookie(ctx.Writer, "refresh_token", "", -1)
	helper.SetCookie(ctx.Writer, "access_token", "", -1)
	responsejson.Success(ctx, nil, "User deleted successfully")
}
