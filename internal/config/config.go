package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	AppEnv    string
	JWTSecret string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	SeedAdminUser string
	SeedAdminPass string
}

func Load() Config {
	// load .env kalau ada (tidak error jika file tidak ada)
	_ = godotenv.Load()

	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		AppEnv:    getEnv("APP_ENV", "dev"),
		JWTSecret: mustEnv("JWT_SECRET"),

		DBHost:     mustEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     mustEnv("DB_NAME"),
		DBUser:     mustEnv("DB_USER"),
		DBPassword: mustEnv("DB_PASSWORD"),
		DBSSLMode:  getEnv("DB_SSLMODE", "require"),

		SeedAdminUser: getEnv("SEED_ADMIN_USER", ""),
		SeedAdminPass: getEnv("SEED_ADMIN_PASS", ""),
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}
