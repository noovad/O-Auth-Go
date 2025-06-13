//go:build wireinject

package api

import (
	"go_auth-project/api/controller"
	"go_auth-project/api/repository"
	"go_auth-project/api/service"
	"go_auth-project/config"

	"github.com/google/wire"
)

func AuthInjector() *controller.AuthController {
	wire.Build(controller.NewAuthController, service.NewUserServiceImpl, service.NewAuthService, repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}
