package service

import (
	"github.com/noovad/go-auth/api/repository"
	"github.com/noovad/go-auth/dto"
	"github.com/noovad/go-auth/helper"
	"github.com/noovad/go-auth/model"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(user dto.CreateUsersRequest) (dto.UserResponse, error)
	FindByEmail(Email string) (dto.UserResponse, error)
	FindByUsername(username string) (dto.UserResponse, error)
	DeleteById(id uuid.UUID) error
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

func (t *UserServiceImpl) CreateUser(user dto.CreateUsersRequest) (dto.UserResponse, error) {
	err := t.Validate.Struct(user)
	if err != nil {
		return dto.UserResponse{}, helper.ErrFailedValidationWrap(err)
	}

	hashPassword, err := helper.HashPassword(user.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	userModel := model.User{
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Password:   hashPassword,
		AvatarType: user.AvatarType,
	}

	createdUser, err := t.UsersRepository.CreateUser(userModel)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Id:         createdUser.Id,
		Name:       createdUser.Name,
		Username:   createdUser.Username,
		Email:      createdUser.Email,
		AvatarType: createdUser.AvatarType,
	}, nil
}

func (t *UserServiceImpl) FindByEmail(Email string) (dto.UserResponse, error) {
	userData, err := t.UsersRepository.FindByEmail(Email)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Id:         userData.Id,
		Name:       userData.Name,
		Username:   userData.Username,
		Email:      userData.Email,
		AvatarType: userData.AvatarType,
	}, nil
}

func (t *UserServiceImpl) FindByUsername(username string) (dto.UserResponse, error) {
	user, err := t.UsersRepository.FindByUsername(username)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Id:         user.Id,
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		AvatarType: user.AvatarType,
		Password:   user.Password,
	}, nil

}

func (t *UserServiceImpl) DeleteById(id uuid.UUID) error {
	return t.UsersRepository.DeleteById(id)
}
