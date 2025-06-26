package repository

import (
	"github.com/noovad/go-auth/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	CreateUser(user model.User) (model.User, error)
	FindByEmail(Email string) (model.User, error)
	FindByUsername(username string) (model.User, error)
	DeleteById(id uuid.UUID) error
	UpdateUser(user model.User) (model.User, error)
}

func NewUsersREpositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func (r *UsersRepositoryImpl) CreateUser(user model.User) (model.User, error) {
	result := r.Db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}
	user.Id = result.Statement.Model.(*model.User).Id
	return user, nil
}

func (t *UsersRepositoryImpl) FindByEmail(Email string) (model.User, error) {
	var user model.User
	result := t.Db.Where("email = ?", Email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (t *UsersRepositoryImpl) FindByUsername(username string) (model.User, error) {
	var user model.User
	result := t.Db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UsersRepositoryImpl) DeleteById(id uuid.UUID) error {
	result := r.Db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UsersRepositoryImpl) UpdateUser(user model.User) (model.User, error) {
	result := r.Db.Model(&model.User{}).
		Where("id = ?", user.Id).
		Update("avatar_type", user.AvatarType)
	if result.Error != nil {
		return user, result.Error
	}
	var updatedUser model.User
	if err := r.Db.First(&updatedUser, "id = ?", user.Id).Error; err != nil {
		return user, err
	}
	return updatedUser, nil
}
