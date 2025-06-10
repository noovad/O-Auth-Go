package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"go_auth-project/config"
	"go_auth-project/data"
	"go_auth-project/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, error)
	CreateTokens(ctx *gin.Context, userId int) error
}

type authService struct {
	usersService UsersService
}

func NewAuthService(usersService UsersService) AuthService {
	return &authService{usersService: usersService}
}

func (s *authService) AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, error) {
	cookieState, err := ctx.Cookie("oauthstate")
	if err != nil || state != cookieState {
		return data.UserResponse{}, helper.ErrInvalidOAuthState
	}
	ctx.SetCookie("oauthstate", "", -1, "/", "", false, true)

	content, err := s.getUserInfoFromGoogle(code)
	if err != nil {
		return data.UserResponse{}, err
	}

	var loginAccount data.GoogleResponse
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		return data.UserResponse{}, err
	}

	user, err := s.usersService.FindByEmail(loginAccount.Email)
	if errors.Is(err, helper.ErrUserNotFound) {
		userId, err := s.usersService.CreateAndReturnID(data.CreateUsersRequest{
			Username: loginAccount.Name,
			Email:    loginAccount.Email,
		})
		if err != nil {
			return data.UserResponse{}, err
		}
		user = data.UserResponse{Id: userId, Email: loginAccount.Email}
	} else if err != nil {
		return data.UserResponse{}, err
	}

	return user, nil
}

func (s *authService) getUserInfoFromGoogle(code string) ([]byte, error) {
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, helper.ErrCodeExchangeFailed(err)
	}

	// Note: This function can be improved by (Suggested by ChatGPT):
	// - Using `http.Client` with a timeout to prevent indefinite hangs.
	// - Sending the access token via `Authorization` header instead of query parameters for better security.
	// - Checking the HTTP response status code before reading the body.
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

func (s *authService) CreateTokens(ctx *gin.Context, userId int) error {
	if err := helper.CreateAccessToken(ctx, userId); err != nil {
		return err
	}
	return helper.CreateRefreshToken(ctx, userId)
}
