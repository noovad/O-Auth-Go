// go:build wireinject
//go:build wireinject
// +build wireinject

package api

import (
	"learn_o_auth-project/api/controller"
	"learn_o_auth-project/api/repository"
	"learn_o_auth-project/api/service"
	"learn_o_auth-project/config"

	"github.com/google/wire"
)

func InitializeUserController() *controller.UsersController {
	wire.Build(controller.NewUsersController, service.NewUsersServiceImpl, repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}

func InitializeAuthController() *controller.UsersAuthController {
	wire.Build(controller.NewUsersAuthController ,service.NewUsersServiceImpl, service.NewAuthService ,repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}
