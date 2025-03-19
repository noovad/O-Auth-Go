package helper

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
func SuccessResponse(ctx *gin.Context, action string, data interface{}) {
	var code int
	var message string

	switch action {
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
		message = "Success"
	}

	ctx.JSON(code, Response{
		Code:   code,
		Status: message,
		Data:   data,
	})
}

func UnauthorizedResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		Code:   http.StatusUnauthorized,
		Status: "Unauthorized",
		Data:   nil,
	})
}

func InternalServerErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code:   http.StatusInternalServerError,
		Status: "Internal Server Error",
		Data:   err.Error(),
	})
}

func BadRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Code:   http.StatusBadRequest,
		Status: "Bad Request",
		Data:   err.Error(),
	})
}

func ForbiddenResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, Response{
		Code:   http.StatusForbidden,
		Status: "Forbidden",
		Data:   message,
	})
}

// func NotFoundResponse(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusNotFound, Response{
// 		Code:   http.StatusNotFound,
// 		Status: "Not Found",
// 		Data:   message,
// 	})
// }

// func ConflictResponse(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusConflict, Response{
// 		Code:   http.StatusConflict,
// 		Status: "Conflict",
// 		Data:   message,
// 	})
// }

// func UnprocessableEntityResponse(ctx *gin.Context, err error) {
// 	ctx.JSON(http.StatusUnprocessableEntity, Response{
// 		Code:   http.StatusUnprocessableEntity,
// 		Status: "Unprocessable Entity",
// 		Data:   err.Error(),
// 	})
// }

// func TooManyRequestsResponse(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusTooManyRequests, Response{
// 		Code:   http.StatusTooManyRequests,
// 		Status: "Too Many Requests",
// 		Data:   message,
// 	})
// }
