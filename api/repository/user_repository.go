package repository

import (
	"errors"
	"go_auth-project/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	SaveAndReturnID(user model.User) (uuid.UUID, error)
	FindByEmail(Email string) (model.User, error)
	FindByUsername(username string) (model.User, error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (r *UsersRepositoryImpl) SaveAndReturnID(user model.User) (uuid.UUID, error) {
	result := r.Db.Create(&user)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return user.Id, nil
}

func (t *UsersRepositoryImpl) FindByEmail(Email string) (model.User, error) {
	var user model.User
	result := t.Db.Where("email = ?", Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, gorm.ErrRecordNotFound
	} else if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (t *UsersRepositoryImpl) FindByUsername(username string) (model.User, error) {
	var user model.User
	result := t.Db.Where("username = ?", username).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, gorm.ErrRecordNotFound
	} else if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
