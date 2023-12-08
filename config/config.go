package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	PromToken string
	PromURL   string
	// Add other configuration fields as needed
}

// LoadEnv loads environment variables and populates the Config struct
func LoadEnv(files ...string) *Config {
	// Load .env file(s)
	for _, file := range files {
		if err := godotenv.Load(file); err != nil {
			log.Printf("Error loading .env file from path %s: %v", file, err)
		}
	}

	// Return a pointer to a Config struct
	// DefaultValues can be set here, if you don't want to deal with .env files
	// Please note, that the default values here shouldn't be used as is
	// They are only here to demonstrate the concept
	return &Config{
		DBHost:     getEnv("DB_HOST", "defaultHost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "mydb"),
		PromToken:  getEnv("PROMETHEUS_TOKEN", "defaultToken"),
		PromURL:    getEnv("PROMETHEUS_URL", "defaultURL"),
	}
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
