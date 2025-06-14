package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port        string
	DataDir     string
	MaxFileSize int64
}

// MaxFileSizeMB returns the maximum file size in megabytes
func (c *Config) MaxFileSizeMB() float64 {
	return float64(c.MaxFileSize) / (1024 * 1024)
}

// Load returns application configuration with defaults
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DataDir:     getEnv("DATA_DIR", "data"),
		MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 100*1024*1024), // 100MB default
	}
	return cfg
}

// getEnv gets environment variable with fallback to default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt64 gets environment variable as int64 with fallback to default
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
