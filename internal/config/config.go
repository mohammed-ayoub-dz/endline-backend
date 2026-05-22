package config

import (
    "log"
    "os"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
    DBSSLMode  string

    AccessSecret  string
    RefreshSecret string
    AccessTTL     time.Duration
    RefreshTTL    time.Duration
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    accessTTL, _ := time.ParseDuration(getEnv("ACCESS_TTL", "15m"))
    refreshTTL, _ := time.ParseDuration(getEnv("REFRESH_TTL", "720h"))

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "postgres"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        AccessSecret:  getEnv("ACCESS_SECRET", "change-me"),
        RefreshSecret: getEnv("REFRESH_SECRET", "change-me-too"),
        AccessTTL:     accessTTL,
        RefreshTTL:    refreshTTL,
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}