package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	JWTSecret     string
	JWTExpireHour int
}

var cfg *Config

func Load() {
	if os.Getenv("APP_NEW") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("no .env file found")
		}
	}

	expireHour, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOUR", "24"))
	cfg = &Config{
		AppEnv:  getEnv("DB_HOST", "localhost"),
		AppPort: getEnv("APP_PORT", "localhost"),

		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "user"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "edo_new"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "secretmakmak"),
		JWTExpireHour: expireHour,
	}
}

func Get() *Config {
	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
