package config

import (
	"fmt"
	"os"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetPort() string {
	port := getEnv("APP_PORT", "3013")
	return fmt.Sprintf(":%s", port)
}

func GetEnv(key string, value string) string {
	return getEnv(key, value)
}
