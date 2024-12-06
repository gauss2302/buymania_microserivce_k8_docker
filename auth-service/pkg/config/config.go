package config

import (
	"os"
	"time"
)

type Config struct {
	// Server settings
	ServerPort string

	// Memcached settings
	MemcachedHost string
	MemcachedPort string

	// Auth settings
	TokenExpiration time.Duration

	// External services
	UserServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		MemcachedHost:   getEnv("MEMCACHED_HOST", "localhost"),
		MemcachedPort:   getEnv("MEMCACHED_PORT", "11211"),
		TokenExpiration: getDurationEnv("TOKEN_EXPIRATION", 24*time.Hour),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://user-service:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
