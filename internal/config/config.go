package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName             string
	Environment         string
	ServerPort          uint16
	DBHost              string
	DBPort              uint16
	DBUser              string
	DBPassword          string
	DBName              string
	DatabaseURL         string
	DBSSLMode           string
	RedisAddr           string
	RedisPassword       string
	JWTSecret           string
	JWTIssuer           string
	JWTAccessDuration   time.Duration
	JWTRefreshDuration  time.Duration
	AfroMessageToken    string
	AfroMessageSenderID string
	AfroMessageURL      string
}

func Load() (*Config, error) {

	cfg := &Config{
		AppName:             getEnv("APP_NAME", "user-service"),
		Environment:         getEnv("APP_ENV", "development"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "user_service"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", ""),
		JWTSecret:           getEnv("JWT_SECRET", "change-me"),
		JWTIssuer:           getEnv("JWT_ISSUER", "user-service"),
		AfroMessageToken:    getEnv("AFROMESSAGE_TOKEN", ""),
		AfroMessageSenderID: getEnv("AFROMESSAGE_SENDER_ID", ""),
		AfroMessageURL:      getEnv("AFROMESSAGE_URL", ""),
	}

	sport := getEnv("SERVER_PORT", "8080")
	port, err := strconv.ParseUint(sport, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	cfg.ServerPort = uint16(port)

	sdbport := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.ParseUint(sdbport, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	cfg.DBPort = uint16(dbPort)

	// FIXED: Removed trailing comma and parentheses issues
	databaseURL := getEnv(
		"DATABASE_URL",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			dbPort,
			cfg.DBName,
			cfg.DBSSLMode,
		),
	)
	cfg.DatabaseURL = databaseURL

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}
	// time
	accessTokenDu := 15 * time.Minute
	refreshTokenDu := 7 * 24 * time.Hour

	cfg.JWTAccessDuration = accessTokenDu
	cfg.JWTRefreshDuration = refreshTokenDu

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
