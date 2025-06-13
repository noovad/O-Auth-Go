package responsejson

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// Create, Read, Update, Delete
func Success(ctx *gin.Context, method string, data interface{}) {
	var code int
	var message string

	switch method {
	case "create":
		code = http.StatusCreated
		message = "Successfully created"
	case "read":
		code = http.StatusOK
		message = "Successfully retrieved data"
	case "update":
		code = http.StatusOK
		message = "Successfully updated"
	case "delete":
		code = http.StatusOK
		message = "Successfully deleted"
	default:
		code = http.StatusOK
		message = method
	}

	ctx.JSON(code, Response{
		Code:   code,
		Status: message,
		Data:   data,
	})
}

func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:   http.StatusUnauthorized,
		Status: "Unauthorized",
		Data:   nil,
	})
}

func InternalServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:   http.StatusInternalServerError,
		Status: "Internal Server Error",
		Data:   err.Error(),
	})
}

func BadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Code:   http.StatusBadRequest,
		Status: "Bad Request",
		Data:   err.Error(),
	})
}

func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:   http.StatusForbidden,
		Status: "Forbidden",
		Data:   message,
	})
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNotFound, Response{
		Code:   http.StatusNotFound,
		Status: "Not Found",
		Data:   message,
	})
}

func Conflict(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, Response{
		Code:   http.StatusConflict,
		Status: "Conflict",
		Data:   message,
	})
}

// func UnprocessableEntity(ctx *gin.Context, err error) {
// 	ctx.JSON(http.StatusUnprocessableEntity, Response{
// 		Code:   http.StatusUnprocessableEntity,
// 		Status: "Unprocessable Entity",
// 		Data:   err.Error(),
// 	})
// }

// func TooManyRequests(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusTooManyRequests, Response{
// 		Code:   http.StatusTooManyRequests,
// 		Status: "Too Many Requests",
// 		Data:   message,
// 	})
// }
