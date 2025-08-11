package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
	DBName      string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "host=localhost user=myuser dbname=meetingschedular port=5432 password=mypassword sslmode=disable"),
		Port:        getEnvOrDefault("PORT", "8080"),
		DBName:      getEnvOrDefault("DBNAME", "meetingschedular"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDSN() string {
	return os.Getenv("DATABASE_URL")
}
