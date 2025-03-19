package service

import (
	"context"
	"io"
	"learn_o_auth-project/config"
	"learn_o_auth-project/helper"
	"net/http"
)

type AuthService interface {
	GetUserInfo(state string, code string) ([]byte, error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) GetUserInfo(state string, code string) ([]byte, error) {
	if state != config.OauthStateString {
		return nil, helper.ErrInvalidOAuthState
	}

	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, helper.ErrCodeExchangeFailed(err)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, helper.ErrFailedGetUserInfo(err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, helper.ErrFailedReadResponseBody(err)
	}

	return contents, nil
}
