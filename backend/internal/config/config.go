package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	AppEnv      string
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	FrontendURL string
	DBSSLMode   string
}

var App Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using system environment variables.")
	}
	sslMode :=
		os.Getenv("DB_SSLMODE")

	if sslMode == "" {
		sslMode = "disable"
	}

	App = Config{
		AppName:     os.Getenv("APP_NAME"),
		AppEnv:      os.Getenv("APP_ENV"),
		AppPort:     os.Getenv("APP_PORT"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSSLMode:   sslMode,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}

	if err := validateConfig(); err != nil {
		log.Fatal(err)
	}
}

func validateConfig() error {
	required := map[string]string{
		"APP_NAME":     App.AppName,
		"APP_ENV":      App.AppEnv,
		"APP_PORT":     App.AppPort,
		"DB_HOST":      App.DBHost,
		"DB_PORT":      App.DBPort,
		"DB_USER":      App.DBUser,
		"DB_NAME":      App.DBName,
		"JWT_SECRET":   App.JWTSecret,
		"FRONTEND_URL": App.FrontendURL,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return errors.New(key + " is required")
		}
	}

	if len(App.JWTSecret) < 16 {
		return errors.New("JWT_SECRET must be at least 16 characters")
	}

	return nil
}
