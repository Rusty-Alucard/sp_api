// Package config provides configuration management for the application.
package config

import "os"

// Config holds the application configuration.
type Config struct {
	Environment string
	Database    DatabaseConfig
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var globalConfig *Config

// Init initializes the global configuration.
func Init(env string) error {
	globalConfig = &Config{
		Environment: env,
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "password"),
			Name:     getEnvOrDefault("DB_NAME", "sp_api"),
		},
	}

	return nil
}

// Get returns the global configuration.
func Get() *Config {
	if globalConfig == nil {
		panic("config not initialized")
	}
	return globalConfig
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
