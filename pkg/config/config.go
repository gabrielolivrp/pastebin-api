package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Test        Environment = "test"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
}

type Config struct {
	Environment Environment
	Port        string
	DB          DatabaseConfig
	Cache       CacheConfig
}

func LoadConfig(filename string) (*Config, error) {
	_ = godotenv.Load(filename)

	envS := os.Getenv("ENV")
	var environment Environment
	switch envS {
	case "production":
		environment = Production
	case "test":
		environment = Test
	default:
		environment = Development
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSL"),
	}

	cache := CacheConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	return &Config{
		Environment: environment,
		Port:        port,
		DB:          db,
		Cache:       cache,
	}, nil
}
