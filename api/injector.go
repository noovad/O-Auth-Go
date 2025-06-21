//go:build wireinject

package api

import (
	"go-auth/api/controller"
	"go-auth/api/repository"
	"go-auth/api/service"
	"go-auth/config"

	"github.com/google/wire"
)

func AuthInjector() *controller.AuthController {
	wire.Build(controller.NewAuthController, service.NewUserServiceImpl, service.NewAuthService, repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}
