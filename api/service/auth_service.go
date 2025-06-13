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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthService interface {
	AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, error)
	AuthenticateWithUsername(ctx *gin.Context, user data.LoginRequest) (uuid.UUID, error)
	CreateTokens(ctx *gin.Context, userId uuid.UUID) error
}

type authService struct {
	usersService UserService
}

func NewAuthService(usersService UserService) AuthService {
	return &authService{usersService: usersService}
}

func (s *authService) AuthenticateWithGoogle(ctx *gin.Context, state string, code string) (data.UserResponse, error) {
	cookieState, err := ctx.Cookie("Oauth-State")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return data.UserResponse{}, helper.ErrOAuthStateNotFound
		}
		return data.UserResponse{}, err
	}

	if state != cookieState {
		return data.UserResponse{}, helper.ErrInvalidOAuthState
	}

	ctx.SetCookie("Oauth-State", "", -1, "/", os.Getenv("FRONTEND_DOMAIN"), false, true)

	content, err := s.getUserInfoFromGoogle(code)
	if err != nil {
		return data.UserResponse{}, err
	}

	var loginAccount data.GoogleResponse
	if err := json.Unmarshal(content, &loginAccount); err != nil {
		return data.UserResponse{}, err
	}

	user, err := s.usersService.FindByEmail(loginAccount.Email)

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			user.Email = loginAccount.Email
			return user, helper.ErrUserNotFound
		}
		return data.UserResponse{}, err
	}

	return user, nil
}

func (s *authService) AuthenticateWithUsername(ctx *gin.Context, user data.LoginRequest) (uuid.UUID, error) {
	existingUser, err := s.usersService.FindByUsername(user.Username)
	if err != nil {
		return uuid.Nil, helper.ErrInvalidCredentials
	}

	if !helper.CheckPasswordHash(user.Password, existingUser.Password) {
		return uuid.Nil, helper.ErrInvalidCredentials
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

func (s *authService) CreateTokens(ctx *gin.Context, userId uuid.UUID) error {
	if err := helper.CreateAccessToken(ctx, userId.String()); err != nil {
		return err
	}
	return helper.CreateRefreshToken(ctx, userId.String())
}
