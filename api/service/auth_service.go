package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"learn_o_auth-project/config"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"
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

	content, err := s.getUserInfoFromGoogle(code)
	if err != nil {
		return data.UserResponse{}, err
	}

	var loginAccount map[string]interface{}
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		return data.UserResponse{}, err
	}

	email, ok := loginAccount["email"].(string)
	if !ok {
		return data.UserResponse{}, helper.ErrFailedToGetEmail
	}

	user, err := s.usersService.FindByEmail(email)
	if errors.Is(err, helper.ErrUserNotFound) {
		userId, err := s.usersService.CreateAndReturnID(data.CreateUsersRequest{
			Username: helper.GetUsernameFromEmail(email),
			Email:    email,
		})
		if err != nil {
			return data.UserResponse{}, err
		}
		user = data.UserResponse{Id: userId, Email: email}
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
