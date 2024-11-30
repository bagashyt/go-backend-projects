package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	DBConfig   string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	publicHost := getEnv("PUBLIC_HOST", "http://localhost")
	port := getEnv("PORT", "8080")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "rahasia123")
	dbAddress := getEnv("DB_HOST", "localhost")
	dbName := getEnv("DB_NAME", "blogs")

	return Config{
		PublicHost: publicHost,
		Port:       port,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBAddress:  fmt.Sprintf("%s:%s", dbAddress),
		DBName:     dbName,
		DBConfig:   getEnv("DB_CFG", fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassword, dbAddress, dbName)),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
