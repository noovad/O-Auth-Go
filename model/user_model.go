package model

import "time"

type Users struct {
	Id        int       `gorm:"primary_key"`
	Username  string    `gorm:"type:varchar(255)"`
	Email     string    `gorm:"type:varchar(255);unique"`
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}
