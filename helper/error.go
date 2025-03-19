package helper

import "errors"

// ✅ Sentinel errors
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden access")
	ErrInternalServer = errors.New("internal server error")
)
