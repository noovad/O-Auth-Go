package model

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}
