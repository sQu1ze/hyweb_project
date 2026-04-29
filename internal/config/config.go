package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	MySQLHost      string
	MySQLPort      string
	MySQLDatabase  string
	MySQLUser      string
	MySQLPassword  string
	JWTSecret      string
	JWTExpireHours int
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "local"),
		MySQLHost:      getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:      getEnv("MYSQL_PORT", "3306"),
		MySQLDatabase:  getEnv("MYSQL_DATABASE", "app_db"),
		MySQLUser:      getEnv("MYSQL_USER", "app_user"),
		MySQLPassword:  getEnv("MYSQL_PASSWORD", "app_pass"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
