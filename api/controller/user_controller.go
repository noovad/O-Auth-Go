package controller

import (
	"errors"
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"

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
	if err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	err = controller.usersService.Create(createUsersRequest)
	if err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, "create", nil)
}

func (controller *UsersController) FindByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	userResponse, err := controller.usersService.FindByEmail(email)

	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			helper.NotFoundResponse(ctx, "User not found")
			return
		}
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, "read", userResponse)
}
