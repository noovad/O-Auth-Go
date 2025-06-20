package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func DatabaseConnection() *gorm.DB {
	once.Do(func() {
		var err error

		var dsn string
		if os.Getenv("ENV") == "production" {
			fmt.Println("Database connection in production mode")
			dsn = os.Getenv("DATABASE_URL")
		} else {
			host := os.Getenv("DBHOST")
			user := os.Getenv("DBUSER")
			password := os.Getenv("DBPASSWORD")
			dbname := os.Getenv("DBNAME")
			port := os.Getenv("DBPORT")

			if host == "" || user == "" || password == "" || dbname == "" || port == "" {
				log.Fatal("One or more required environment variables are not set")
			}
			fmt.Println("Database connection in development mode")
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				host, user, password, dbname, port)
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("Failed to get raw database object:", err)
		}
		log.Println("Database connected successfully")

		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

	})

	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get DB object for closing:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}
}
