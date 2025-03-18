package service

import (
	"learn_o_auth-project/api/repository"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/model"

	"github.com/go-playground/validator/v10"
)

type UsersService interface {
	Create(Users data.CreateUsersRequest)
	FindByEmail(Email string) data.UsersResponse
}

func NewUsersServiceImpl(userRepository repository.UsersRepository, validate *validator.Validate) UsersService {
	return &UsersServiceImpl{
		UsersRepository: userRepository,
		Validate:        validate,
	}
}

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
	Validate        *validator.Validate
}

func (t *UsersServiceImpl) Create(Users data.CreateUsersRequest) {
	err := t.Validate.Struct(Users)
	helper.ErrorPanic(err)
	userModel := model.Users{
		Username: Users.Username,
		Email:    Users.Email,
	}
	t.UsersRepository.Save(userModel)
}

func (t *UsersServiceImpl) FindByEmail(Email string) data.UsersResponse {
	userData, err := t.UsersRepository.FindByEmail(Email)
	helper.ErrorPanic(err)

	userResponse := data.UsersResponse{
		Id:       userData.Id,
		Username: userData.Username,
		Email:    userData.Email,
	}
	return userResponse
}
