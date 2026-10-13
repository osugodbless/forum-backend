package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type serverConfig struct {
	host string
	port int
	env  string
}

type DBConfig struct {
	URL string
}

// appConfig combines all configurations
type appConfig struct {
	server serverConfig
	DB     DBConfig
}

func LoadConfig() (*appConfig, error) {
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		return nil, err
	}

	dbURL := getEnv("DB_URL", "")
	if dbURL == "" {
		return nil, errors.New("DB_URL is required")
	}

	return &appConfig{
		server: serverConfig{
			host: getEnv("APP_HOST", "localhost"),
			port: getEnvInt("APP_PORT", 8080),
			env:  getEnv("APP_ENVIRONMENT", "development"),
		},
		DB: DBConfig{
			URL: dbURL,
		},
	}, nil

}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
