package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username   string    `gorm:"type:varchar(255)"`
	Name       string    `gorm:"type:varchar(255)"`
	Email      string    `gorm:"type:varchar(255);unique"`
	AvatarType string    `gorm:"type:varchar(255)"`
	Password   string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"type:timestamp"`
	UpdatedAt  time.Time `gorm:"type:timestamp"`
}
