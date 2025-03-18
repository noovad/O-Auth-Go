package controller

import (
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	usersService service.UsersService
}

func NewUsersController(service service.UsersService) *UsersController {
	return &UsersController{
		usersService: service,
	}
}

func (controller *UsersController) Create(ctx *gin.Context) {
	createUsersRequest := data.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)

	controller.usersService.Create(createUsersRequest)
	webResponse := data.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

func (controller *UsersController) FindByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	userResponse := controller.usersService.FindByEmail(email)

	webResponse := data.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   userResponse,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}
