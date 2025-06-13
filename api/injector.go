// go:build wireinject
//go:build wireinject
// +build wireinject

package api

import (
	"go_auth-project/api/controller"
	"go_auth-project/api/repository"
	"go_auth-project/api/service"
	"go_auth-project/config"

	"github.com/google/wire"
)

func InitializeAuthController() *controller.AuthController {
	wire.Build(controller.NewAuthController, service.NewUsersServiceImpl, service.NewAuthService, repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}
