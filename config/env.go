package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser                 string
	DBPassword             string
	DBAdress               string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
	MigrateSource          string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "root"),
		DBAdress:               fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTSecret:              getEnv("JWT_SECRET", "super-secret-key"), // should be at least 32 characters long
		MigrateSource:          getEnv("MIGRATIONS_PATH", "file://migrations"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7), // 7 days
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return intValue
	}

	return fallback
}
