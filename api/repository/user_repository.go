package repository

import (
	"errors"
	"learn_o_auth-project/helper"
	"learn_o_auth-project/model"

	"gorm.io/gorm"
)

type UsersRepository interface {
	Save(Users model.Users)
	FindByEmail(Email string) (Users model.Users, err error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (t *UsersRepositoryImpl) FindByEmail(Email string) (Users model.Users, err error) {
	var user model.Users
	result := t.Db.Where("email = ?", Email).First(&user)
	if result != nil {
		return user, nil
	} else {
		return user, errors.New("user is not found")
	}
}

func (t *UsersRepositoryImpl) Save(Users model.Users) {
	result := t.Db.Create(&Users)
	helper.ErrorPanic(result.Error)
}
