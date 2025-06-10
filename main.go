package main

import (
	"learn_o_auth-project/config"
	"learn_o_auth-project/model"
	"learn_o_auth-project/router"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.DatabaseConnection()
	if db == nil {
		log.Fatal("Database connection failed")
	}

	if err := model.Migration(db); err != nil {
		log.Fatal("Database migration failed:", err)
	}

	config.InitOAuth()

	r := router.SetupRouter()

	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server running on port", os.Getenv("PORT"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
