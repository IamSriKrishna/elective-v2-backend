package config

import (
	"os"
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
}

type ServerConfig struct {
	Port string
}

func LoadConfig() *Config {
	return &Config{
		Database: DataBaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			DBName:   getEnv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT"),
		},
	}
}

func getEnv(key string) string {

	return os.Getenv(key)
}
