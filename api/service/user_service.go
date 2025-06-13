package service

import (
	"go_auth-project/api/repository"
	"go_auth-project/data"
	"go_auth-project/helper"
	"go_auth-project/model"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserService interface {
	CreateAndReturnID(user data.CreateUsersRequest) (uuid.UUID, error)
	FindByEmail(Email string) (data.UserResponse, error)
	FindByUsername(username string) (data.UserResponse, error)
}

func NewUserServiceImpl(userRepository repository.UsersRepository, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UsersRepository: userRepository,
		Validate:        validate,
	}
}

type UserServiceImpl struct {
	UsersRepository repository.UsersRepository
	Validate        *validator.Validate
}

func (t *UserServiceImpl) CreateAndReturnID(user data.CreateUsersRequest) (uuid.UUID, error) {
	err := t.Validate.Struct(user)
	if err != nil {
		return uuid.Nil, helper.ErrFailedValidationWrap(err)
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return uuid.Nil, err
	}

	userModel := model.User{
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Password:   hashPassword,
		AvatarType: user.AvatarType,
	}

	return t.UsersRepository.SaveAndReturnID(userModel)
}

func (t *UserServiceImpl) FindByEmail(Email string) (data.UserResponse, error) {
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

func (t *UserServiceImpl) FindByUsername(username string) (data.UserResponse, error) {
	userData, err := t.UsersRepository.FindByUsername(username)
	if err != nil {
		return data.UserResponse{}, err
	}

	return data.UserResponse{
		Id:       userData.Id,
		Username: userData.Username,
		Email:    userData.Email,
		Password: userData.Password,
	}, nil
}
