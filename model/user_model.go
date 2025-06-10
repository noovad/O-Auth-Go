package model

import "time"

type Users struct {
	Id         string    `gorm:"primary_key"`
	Username   string    `gorm:"type:varchar(255)"`
	Name       string    `gorm:"type:varchar(255)"`
	Email      string    `gorm:"type:varchar(255);unique"`
	AvatarType string    `gorm:"type:varchar(255)"`
	Password   string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"type:timestamp"`
	UpdatedAt  time.Time `gorm:"type:timestamp"`
}
