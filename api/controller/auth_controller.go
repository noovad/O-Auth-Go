package controller

import (
	"errors"
	"net/http"

	"github.com/noovad/go-auth/api/service"
	"github.com/noovad/go-auth/config"
	"github.com/noovad/go-auth/dto"
	"github.com/noovad/go-auth/helper"
	"github.com/noovad/go-auth/helper/responsejson"

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
	accessTokenAge := helper.MustParseDurationEnv("REFRESH_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, int(accessTokenAge))

	accessToken := helper.CreateAccessToken(user)
	refreshTokenAge := helper.MustParseDurationEnv("ACCESS_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "access_token", accessToken, int(refreshTokenAge))

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
		if errors.Is(err, helper.ErrUsernameNotFound) {
			responsejson.NotFound(ctx, "Username not found")
			return
		}
		if errors.Is(err, helper.ErrWrongPassword) {
			responsejson.Unauthorized(ctx, "Wrong password")
			return
		}
		if errors.Is(err, helper.ErrInvalidCredentials) {
			responsejson.Unauthorized(ctx, "Invalid username or password")
			return
		}
		responsejson.InternalServerError(ctx, err, "Failed to authenticate user")
		return
	}

	refreshToken := helper.CreateRefreshToken(userResponse)
	accessTokenAge := helper.MustParseDurationEnv("REFRESH_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, int(accessTokenAge))

	accessToken := helper.CreateAccessToken(userResponse)
	refreshTokenAge := helper.MustParseDurationEnv("ACCESS_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "access_token", accessToken, int(refreshTokenAge))

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
			ctx.Redirect(http.StatusPermanentRedirect, helper.MustGetenv("FRONTEND_BASE_URL")+"/login?error=Invalid credentials, please login again")
			return
		}
		if errors.Is(err, helper.ErrInvalidOAuthState) {
			ctx.Redirect(http.StatusPermanentRedirect, helper.MustGetenv("FRONTEND_BASE_URL")+"/login?error=Invalid credentials, please login again")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.SetCookie("Signed-token", helper.CreateSignedToken(user.Email), 60*1, "/", helper.MustGetenv("BACKEND_DOMAIN"), true, true)
			ctx.Redirect(http.StatusPermanentRedirect, helper.MustGetenv("FRONTEND_BASE_URL")+"/sign-up?email="+user.Email)
			return
		}

		ctx.Redirect(http.StatusPermanentRedirect, helper.MustGetenv("FRONTEND_BASE_URL")+"/login?error=Failed to authenticate with Google. Please try again later.")
		return
	}

	refreshToken := helper.CreateRefreshToken(user)
	accessTokenAge := helper.MustParseDurationEnv("REFRESH_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "refresh_token", refreshToken, int(accessTokenAge))

	accessToken := helper.CreateAccessToken(user)
	refreshTokenAge := helper.MustParseDurationEnv("ACCESS_TOKEN_AGE")
	helper.SetCookie(ctx.Writer, "access_token", accessToken, int(refreshTokenAge))

	redirectURL := helper.MustGetenv("FRONTEND_BASE_URL") + "/"
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (c *AuthController) HandleDeleteAccount(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
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

func (c *AuthController) HandleUpdateAvatar(ctx *gin.Context) {
	var updateAvatarReq dto.UpdateAvatarUserRequest
	if err := ctx.ShouldBindJSON(&updateAvatarReq); err != nil {
		responsejson.BadRequest(ctx, err, "Invalid update avatar request")
		return
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		responsejson.InternalServerError(ctx, nil, "User ID not found in context")
		return
	}
	uid, ok := userId.(uuid.UUID)
	if !ok {
		responsejson.InternalServerError(ctx, nil, "Invalid user ID type")
		return
	}

	updateAvatarReq.Id = uid
	userResponse, err := c.userService.UpdateAvatar(updateAvatarReq)
	if err != nil {
		responsejson.InternalServerError(ctx, err, "Failed to update avatar")
		return
	}

	responsejson.Success(ctx, gin.H{
		"user": userResponse,
	}, "Avatar updated successfully")
}
