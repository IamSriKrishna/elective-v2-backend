package config

import (
    "os"
    "time"
)

type Config struct {
    Database DataBaseConfig
    JWT      JWTConfig
    Server   ServerConfig
}

type DataBaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

type ServerConfig struct {
    Port string
}

func LoadConfig() *Config {
    return &Config{
        Database: DataBaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "3000"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "2004"),
            DBName:   getEnv("DB_NAME", "course_booking"),
        },
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", "your-secret-key"),
            Expiration: 24 * time.Hour,
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}