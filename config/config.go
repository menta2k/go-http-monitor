package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port              int
	DBPath            string
	DefaultInterval   int
	HTTPClientTimeout int
	JWTSecret         string
	JWTTokenTTLHours  int
	AdminUsername     string
	AdminPassword     string
	TSDBPath         string
	SMTPHost         string
	SMTPPort         int
	SMTPFrom         string
	SMTPUsername     string
	SMTPPassword     string
}

func Load() Config {
	return Config{
		Port:              envInt("PORT", 8080),
		DBPath:            envStr("DB_PATH", "./monitor.db"),
		DefaultInterval:   envInt("DEFAULT_INTERVAL", 60),
		HTTPClientTimeout: envInt("HTTP_CLIENT_TIMEOUT", 30),
		JWTSecret:         envStr("JWT_SECRET", ""),
		JWTTokenTTLHours:  envInt("JWT_TOKEN_TTL_HOURS", 24),
		AdminUsername:     envStr("ADMIN_USERNAME", "admin"),
		AdminPassword:     envStr("ADMIN_PASSWORD", ""),
		TSDBPath:          envStr("TSDB_PATH", "./tsdb-data"),
		SMTPHost:          envStr("SMTP_HOST", ""),
		SMTPPort:          envInt("SMTP_PORT", 587),
		SMTPFrom:          envStr("SMTP_FROM", ""),
		SMTPUsername:      envStr("SMTP_USERNAME", ""),
		SMTPPassword:      envStr("SMTP_PASSWORD", ""),
	}
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
