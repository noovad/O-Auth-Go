package service

import (
	"context"
	"encoding/json"
	"errors"
	"go_auth-project/config"
	"go_auth-project/data"
	"go_auth-project/helper"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, string, error)
	AuthenticateWithPassword(ctx *gin.Context, user data.LoginRequest) (string, error)
	CreateTokens(ctx *gin.Context, userId string) error
}

type authService struct {
	usersService UsersService
}

func NewAuthService(usersService UsersService) AuthService {
	return &authService{usersService: usersService}
}

func (s *authService) AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, string, error) {
	cookieState, err := ctx.Cookie("Oauth-State") // Retrieve the state from the cookie
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return data.UserResponse{}, "", helper.ErrOAuthStateNotFound
		}
		return data.UserResponse{}, "", err
	}

	parts := strings.SplitN(state, "|", 2)
	if len(parts) != 2 {
		return data.UserResponse{}, "", helper.ErrInvalidOAuthState
	}

	parsedState, action := parts[0], parts[1]

	if parsedState != strings.Split(cookieState, "|")[0] {
		return data.UserResponse{}, "", helper.ErrInvalidOAuthState
	}

	ctx.SetCookie("Oauth-State", "", -1, "/", os.Getenv("FRONTEND_DOMAIN"), false, true)

	content, err := s.getUserInfoFromGoogle(code)
	if err != nil {
		return data.UserResponse{}, "", err
	}

	var loginAccount data.GoogleResponse
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		return data.UserResponse{}, "", err
	}

	user, err := s.usersService.FindByEmail(loginAccount.Email)

	switch action {
	case "login":
		if errors.Is(err, helper.ErrUserNotFound) {
			return data.UserResponse{}, action, helper.ErrUserNotFound
		}
		if err != nil {
			return data.UserResponse{}, action, err
		}

	case "signup":
		if errors.Is(err, helper.ErrUserNotFound) {
			return data.UserResponse{
				Email: loginAccount.Email,
			}, action, helper.ErrUserNotFound
		}
		if err != nil {
			return data.UserResponse{}, action, err
		}
	}

	return user, action, nil
}

func (s *authService) AuthenticateWithPassword(ctx *gin.Context, user data.LoginRequest) (string, error) {
	existingUser, err := s.usersService.FindByEmail(user.Email)
	if err != nil {
		return "", helper.ErrInvalidCredentials
	}

	if !helper.CheckPasswordHash(user.Password, existingUser.Password) {
		return "", helper.ErrInvalidCredentials
	}

	return existingUser.Id, nil
}

func (s *authService) getUserInfoFromGoogle(code string) ([]byte, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, helper.ErrCodeExchangeFailed(err)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, helper.ErrFailedGetUserInfo(err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, helper.ErrFailedReadResponseBody(err)
	}

	return content, nil
}

func (s *authService) CreateTokens(ctx *gin.Context, userId string) error {
	if err := helper.CreateAccessToken(ctx, userId); err != nil {
		return err
	}
	return helper.CreateRefreshToken(ctx, userId)
}
