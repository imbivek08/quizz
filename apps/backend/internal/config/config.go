package config

import "os"

type Config struct {
	ServerAddress string
	Environment   string
}

func Load() *Config {
	return &Config{
		ServerAddress: getEnv("PORT", ":8080"),
		Environment:   getEnv("ENVIRONMENT", "dev"),
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
