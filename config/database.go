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
			if os.Getenv("DATABASE_URL") == "" {
				log.Fatal("Environment variable DATABASE_URL is not set or empty")
			}

			dsn = os.Getenv("DATABASE_URL")
			fmt.Println("Database connection in production mode")
		} else {
			keys := []string{"DBHOST", "DBUSER", "DBPASSWORD", "DBNAME", "DBPORT"}
			for _, key := range keys {
				val := os.Getenv(key)
				if val == "" {
					log.Fatalf("Environment variable %s is not set or empty", key)
				}
			}

			host := os.Getenv("DBHOST")
			user := os.Getenv("DBUSER")
			password := os.Getenv("DBPASSWORD")
			dbname := os.Getenv("DBNAME")
			port := os.Getenv("DBPORT")

			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				host, user, password, dbname, port)
			fmt.Println("Database connection in development mode")
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
