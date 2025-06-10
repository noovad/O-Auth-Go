package repository

import (
	"errors"
	"go_auth-project/helper"
	"go_auth-project/model"

	"gorm.io/gorm"
)

type UsersRepository interface {
	SaveAndReturnID(user model.Users) (int, error)
	FindByEmail(Email string) (model.Users, error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (r *UsersRepositoryImpl) SaveAndReturnID(user model.Users) (int, error) {
	result := r.Db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.Id, nil
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
