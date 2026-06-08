package config

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Addr        string
	LogLevel    slog.Level
	DatabaseURL string
	JWTSecret   string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3UseSSL    bool
}

func Load() *Config {
	return &Config{
		Addr:        env("ADDR", ":8080"),
		LogLevel:    parseLogLevel(env("LOG_LEVEL", "info")),
		DatabaseURL: env("DATABASE_URL", "postgres://cubenet:cubenet@postgres:5432/cubenet?sslmode=disable"),
		JWTSecret:   env("JWT_SECRET", "dev-secret-change-in-production"),
		S3Endpoint:  env("S3_ENDPOINT", "minio:9000"),
		S3AccessKey: env("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: env("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    env("S3_BUCKET", "cubenet"),
		S3UseSSL:    envBool("S3_USE_SSL", false),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v == "true" || v == "1"
}

func envDur(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
