package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// server
	Port      string // HTTP gateway port
	GRPCPort  string // gRPC server port
	LogLevel  string
	LogFormat string

	// db
	Database DatabaseConfig

	// auth and jwt
	JWTSecret   string
	JWTTokenTTL time.Duration
}

type DatabaseConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	AuthSource string
}

func Load(envFiles ...string) (*Config, error) {
	files := envFiles
	if len(files) == 0 {
		files = []string{".env"}
	}
	_ = godotenv.Load(files...) // non-fatal; env vars set externally take priority

	cfg := &Config{
		Port:      getEnv("PORT", "5001"),
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		Database: DatabaseConfig{
			Host:       getEnv("MONGO_DB_HOST", "localhost"),
			Port:       getEnvInt("MONGO_DB_PORT", 27017),
			User:       getEnv("MONGO_DB_USERNAME", "root"),
			Password:   getEnv("MONGO_DB_PASSWORD", "root_pass"),
			Name:       getEnv("MONGO_DB", "rbac-db"),
			AuthSource: getEnv("MONGO_DB_AUTH_SOURCE", "admin"),
		},
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production-please-set-32b"),
	}

	// Parse JWT TTL
	ttlStr := getEnv("JWT_TTL", "24h")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL %q: %w", ttlStr, err)
	}
	cfg.JWTTokenTTL = ttl

	return cfg, nil
}

// helpers
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}
