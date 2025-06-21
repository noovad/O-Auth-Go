//go:build wireinject

package api

import (
	"github.com/noovad/go-auth/api/controller"
	"github.com/noovad/go-auth/api/repository"
	"github.com/noovad/go-auth/api/service"
	"github.com/noovad/go-auth/config"

	"github.com/google/wire"
)

func AuthInjector() *controller.AuthController {
	wire.Build(controller.NewAuthController, service.NewUserServiceImpl, service.NewAuthService, repository.NewUsersREpositoryImpl, config.DatabaseConnection, config.NewValidator)
	return nil
}
