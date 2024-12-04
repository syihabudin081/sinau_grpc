package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}
}

// GetEnv retrieves an environment variable or a default value if not found
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
