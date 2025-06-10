package service

import (
	"go_auth-project/api/repository"
	"go_auth-project/data"
	"go_auth-project/helper"
	"go_auth-project/model"

	"github.com/go-playground/validator/v10"
)

type UsersService interface {
	CreateAndReturnID(user data.CreateUsersRequest) (string, error)
	FindByEmail(Email string) (data.UserResponse, error)
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

func (t *UsersServiceImpl) CreateAndReturnID(user data.CreateUsersRequest) (string, error) {
	err := t.Validate.Struct(user)
	if err != nil {
		return "", helper.ErrFailedValidationWrap(err)
	}

	userModel := model.Users{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
		AvatarType: user.AvatarType,
	}

	return t.UsersRepository.SaveAndReturnID(userModel)
}

func (t *UsersServiceImpl) FindByEmail(Email string) (data.UserResponse, error) {
	userData, err := t.UsersRepository.FindByEmail(Email)
	if err != nil {
		return data.UserResponse{}, err
	}

	return data.UserResponse{
		Id:       userData.Id,
		Username: userData.Username,
		Email:    userData.Email,
	}, nil
}
