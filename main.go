package main

import (
	"go_auth-project/config"
	"go_auth-project/model"
	"go_auth-project/router"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load it, skipping...")
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server running on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
