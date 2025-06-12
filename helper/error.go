package helper

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidOAuthState  = errors.New("invalid oauth state")
	ErrFailedValidation   = errors.New("validation failed")
	ErrOAuthStateNotFound = errors.New("oauth state not found")

	ErrCodeExchangeFailed = func(err error) error {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}
	ErrFailedGetUserInfo = func(err error) error {
		return fmt.Errorf("failed getting user info: %s", err.Error())
	}
	ErrFailedReadResponseBody = func(err error) error {
		return fmt.Errorf("failed reading response body: %s", err.Error())
	}
	ErrFailedValidationWrap = func(err error) error {
		return fmt.Errorf("%w: %v", ErrFailedValidation, err)
	}
)
