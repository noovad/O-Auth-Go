package responsejson

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func Success(ctx *gin.Context, data any, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func Created(ctx *gin.Context, data any, message ...string) {
	msg := "Resource created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func Unauthorized(ctx *gin.Context, message ...string) {
	msg := "Unauthorized"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: msg,
	})
}

func InternalServerError(ctx *gin.Context, err error, message ...string) {
	msg := "Internal Server Error"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: msg,
		Errors:  err.Error(),
	})
}

func BadRequest(ctx *gin.Context, err error, message ...string) {
	msg := "Bad Request"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: msg,
		Errors:  err.Error(),
	})
}

func Forbidden(ctx *gin.Context, message ...string) {
	msg := "Forbidden"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: msg,
	})
}

func NotFound(ctx *gin.Context, message ...string) {
	msg := "Not Found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: msg,
	})
}

func Conflict(ctx *gin.Context, err error, message ...string) {
	msg := "Conflict"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	ctx.JSON(http.StatusConflict, Response{
		Success: false,
		Message: msg,
		Errors:  err.Error(),
	})
}
