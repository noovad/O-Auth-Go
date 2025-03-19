package repository

import (
	"errors"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/model"

	"gorm.io/gorm"
)

type UsersRepository interface {
	Save(Users model.Users) error
	FindByEmail(Email string) (model.Users, error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (t *UsersRepositoryImpl) Save(Users model.Users) error {
	result := t.Db.Create(&Users)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *UsersRepositoryImpl) FindByEmail(Email string) (model.Users, error) {
	var user model.Users
	result := t.Db.Where("email = ?", Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, helper.ErrUserNotFound
	} else if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
