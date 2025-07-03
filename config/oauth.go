package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config

func InitOAuth() {
	keys := []string{"GOOGLE_REDIRECT_URL", "GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET"}
	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			log.Fatalf("Environment variable %s is not set or empty", key)
		}
	}
	
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
