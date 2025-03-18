package model

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.Table("users").AutoMigrate(&Users{})
}
