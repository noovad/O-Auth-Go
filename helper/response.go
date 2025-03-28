package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct{}

func (Response) Success(ctx *gin.Context, action string, data interface{}) {
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

	ctx.JSON(code, gin.H{
		"code":   code,
		"status": message,
		"data":   data,
	})
}

func (Response) Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"code":   http.StatusUnauthorized,
		"status": "Unauthorized",
		"data":   nil,
	})
}

func (Response) InternalServerError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":   http.StatusInternalServerError,
		"status": "Internal Server Error",
		"data":   err.Error(),
	})
}

func (Response) Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"code":   http.StatusForbidden,
		"status": "Forbidden",
		"data":   message,
	})
}

// func (Response) BadRequest(ctx *gin.Context, err error) {
// 	ctx.JSON(http.StatusBadRequest, gin.H{
// 		"code":   http.StatusBadRequest,
// 		"status": "Bad Request",
// 		"data":   err.Error(),
// 	})
// }

// func (Response) NotFound(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusNotFound, gin.H{
// 		"code":   http.StatusNotFound,
// 		"status": "Not Found",
// 		"data":   message,
// 	})
// }

// func (Response) Conflict(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusConflict, gin.H{
// 		"code":   http.StatusConflict,
// 		"status": "Conflict",
// 		"data":   message,
// 	})
// }

// func (Response) UnprocessableEntity(ctx *gin.Context, err error) {
// 	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
// 		"code":   http.StatusUnprocessableEntity,
// 		"status": "Unprocessable Entity",
// 		"data":   err.Error(),
// 	})
// }

// func (Response) TooManyRequests(ctx *gin.Context, message string) {
// 	ctx.JSON(http.StatusTooManyRequests, gin.H{
// 		"code":   http.StatusTooManyRequests,
// 		"status": "Too Many Requests",
// 		"data":   message,
// 	})
// }
