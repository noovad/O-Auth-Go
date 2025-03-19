package service

import (
	"learn_o_auth-project/api/repository"
	"learn_o_auth-project/data"
	"learn_o_auth-project/model"

	"github.com/go-playground/validator/v10"
)

type UsersService interface {
	Create(Users data.CreateUsersRequest) error
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

func (t *UsersServiceImpl) Create(Users data.CreateUsersRequest) error {
	err := t.Validate.Struct(Users)
	if err != nil {
		return err
	}

	userModel := model.Users{
		Username: Users.Username,
		Email:    Users.Email,
	}

	err = t.UsersRepository.Save(userModel)
	if err != nil {
		return err
	}
	return nil
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
