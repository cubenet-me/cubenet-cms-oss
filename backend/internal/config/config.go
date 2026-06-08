package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr        string
	LogLevel    slog.Level
	DatabaseURL string
	DBName      string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Addr:        env("ADDR", ":8080"),
		LogLevel:    parseLogLevel(env("LOG_LEVEL", "info")),
		DatabaseURL: env("DATABASE_URL", "postgres://cubenet:cubenet@localhost:5432/cubenet?sslmode=disable"),
		DBName:      env("DB_NAME", "cubenet"),
		JWTSecret:   env("JWT_SECRET", "dev-secret-change-in-production"),
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
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func envInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return n
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
