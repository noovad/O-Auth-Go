package repository

import (
	"errors"
	"go_auth-project/helper"
	"go_auth-project/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	SaveAndReturnID(user model.Users) (uuid.UUID, error)
	FindByEmail(Email string) (model.Users, error)
	FindByUsername(username string) (model.Users, error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (r *UsersRepositoryImpl) SaveAndReturnID(user model.Users) (uuid.UUID, error) {
	result := r.Db.Create(&user)
	if result.Error != nil {
		return uuid.Nil, result.Error
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

func (t *UsersRepositoryImpl) FindByUsername(username string) (model.Users, error) {
	var user model.Users
	result := t.Db.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, helper.ErrUserNotFound
	} else if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
