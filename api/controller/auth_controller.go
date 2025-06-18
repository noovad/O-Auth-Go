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
	"time"

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
	ctx.SetCookie("Oauth-State", state, 60, "/", os.Getenv("BACKEND_DOMAIN"), true, true)

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
			responsejson.Unauthorized(ctx, "Invalid signed token")
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

	accessToken := helper.CreateAccessToken(user.Id.String())
	refreshToken := helper.CreateRefreshToken(user.Id.String())

	responsejson.Success(ctx, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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

	accessToken := helper.CreateAccessToken(userResponse.Id.String())
	refreshToken := helper.CreateRefreshToken(userResponse.Id.String())

	responsejson.Success(ctx, gin.H{
		"user":          userResponse,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, "Successfully logged in")
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
	oneTimeCode := uuid.NewString()
	accessToken := helper.CreateAccessToken(user.Id.String())
	refreshToken := helper.CreateRefreshToken(user.Id.String())

	helper.SetMemory(oneTimeCode, helper.TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, time.Minute*2)

	redirectURL := os.Getenv("FRONTEND_BASE_URL") + "/callback?code=" + oneTimeCode
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (c *AuthController) ExchangeCode(ctx *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	
	if err := ctx.BindJSON(&req); err != nil || req.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tokenData, ok := helper.Getmemory(req.Code)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Code expired or invalid"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken":  tokenData.AccessToken,
		"refreshToken": tokenData.RefreshToken,
		"user":         tokenData.User,
	})
}

func (c *AuthController) HandleRefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("refreshToken")

	userId, valid := helper.ValidateRefreshToken(refreshToken)
	if !valid {
		responsejson.Unauthorized(ctx, "Invalid refresh token")
		return
	}

	exists, err := helper.UserExistsInDatabase(userId)
	if err != nil {
		responsejson.InternalServerError(ctx, err, "Failed to check user existence")
		return
	}

	if !exists {
		responsejson.Unauthorized(ctx, "User not found")
		return
	}

	accessToken := helper.CreateAccessToken(userId)

	responsejson.Success(ctx, gin.H{
		"access_token": accessToken,
	}, "Token refreshed successfully")
}

func (c *AuthController) HandleDeleteAccount(ctx *gin.Context) {
	id, valid := helper.ValidateAccessToken(helper.AccessTokenFromHeader(ctx))
	if !valid {
		responsejson.Unauthorized(ctx, "Invalid access token")
		return
	}

	UUID, err := helper.StringToUUID(id)
	if err != nil {
		responsejson.BadRequest(ctx, err, "Invalid user ID format")
		return
	}

	err = c.userService.DeleteById(UUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responsejson.NotFound(ctx, "User not found")
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to delete user")
		return
	}

	responsejson.Success(ctx, nil, "User deleted successfully")
}
