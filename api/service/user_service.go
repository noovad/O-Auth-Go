package service

import (
	"learn_o_auth-project/api/repository"
	"learn_o_auth-project/data"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/model"

	"github.com/go-playground/validator/v10"
)

type UsersService interface {
	CreateAndReturnID(user data.CreateUsersRequest) (int, error)
	FindByEmail(Email string) (data.UsersResponse, error)
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

func (t *UsersServiceImpl) CreateAndReturnID(user data.CreateUsersRequest) (int, error) {
	err := t.Validate.Struct(user)
	if err != nil {
		return 0, helper.ErrFailedValidationWrap(err)
	}

	userModel := model.Users{
		Username: user.Username,
		Email:    user.Email,
	}

	return t.UsersRepository.SaveAndReturnID(userModel)
}

func (t *UsersServiceImpl) FindByEmail(Email string) (data.UsersResponse, error) {
	userData, err := t.UsersRepository.FindByEmail(Email)
	if err != nil {
		return data.UsersResponse{}, err
	}

	return data.UsersResponse{
		Id:       userData.Id,
		Username: userData.Username,
		Email:    userData.Email,
	}, nil
}
