package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	RateLimit   int
	RateBurst   int
	LogLevel    string
	LogFormat   string
	CORSOrigins string
}

var AppConfig *Config

func LoadConfig() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "postcomments"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		LogFormat:   getEnv("LOG_FORMAT", "json"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}
	cfg.RateLimit, _ = strconv.Atoi(getEnv("RATE_LIMIT", "5"))
	cfg.RateBurst, _ = strconv.Atoi(getEnv("RATE_BURST", "10"))
	AppConfig = cfg
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
