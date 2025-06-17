package helper

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidOAuthState   = errors.New("invalid oauth state")
	ErrOAuthStateNotFound  = errors.New("oauth state not found")

	ErrCodeExchangeFailed = func(err error) error {
		return fmt.Errorf("code exchange failed: %s", err.Error())
	}

	ErrFailedValidation     = errors.New("validation failed")
	ErrFailedValidationWrap = func(err error) error {
		return fmt.Errorf("%w: %v", ErrFailedValidation, err)
	}
)
