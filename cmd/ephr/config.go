package main

import (
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/middleware"
	"log"
	"os"
	"strconv"
	"strings"
)

type config struct {
	port  int32
	env   string
	db    database.Config
	auth  middleware.AuthenticateConfig
	cors  middleware.CorsConfig
	limit middleware.RateLimitConfig
}

func parseConfig() config {
	return config{
		port: getEnvInt32Value("PORT", 4000),
		env:  getEnvStrValue("ENV", "dev"),
		db: database.Config{
			DSN:                   getEnvStrValue("DB_URL", ""),
			MaxConns:              getEnvInt32Value("DB_MAX_CONNS", 25),
			MinConns:              getEnvInt32Value("DB_MIN_CONNS", 25),
			MaxConnLifetime:       getEnvStrValue("DB_MAX_CONN_LIFETIME", "15m"),
			MaxConnLifetimeJitter: getEnvStrValue("DB_MAX_CONN_LIFETIME_JITTER", "15m"),
			MaxConnIdleTime:       getEnvStrValue("DB_MAX_CONN_IDLE_TIME", "60m"),
		},
		auth: middleware.AuthenticateConfig{
			JWKSURL: getEnvStrValue("AUTH_JWKS_URL", ""),
			APIURL:  getEnvStrValue("AUTH_API_URL", ""),
		},
		cors: middleware.CorsConfig{
			TrustedOrigins: strings.Fields(getEnvStrValue("CORS_TRUSTED_ORIGINS", "")),
		},
		limit: middleware.RateLimitConfig{
			Enabled: getEnvBoolValue("RATE_LIMIT_ENABLED", true),
			RPS:     getEnvFloat64Value("RATE_LIMIT_RPS", 2),
			Burst:   int(getEnvInt32Value("RATE_LIMIT_BURST", 2)),
		},
	}
}

func getEnvStrValue(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolValue(key string, defaultValue bool) bool {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		switch value {
		case "true", "True", "1":
			return true
		default:
			return false
		}
	}
	return defaultValue
}

func getEnvInt32Value(key string, defaultValue int32) int32 {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}

		return int32(i)
	}
	return defaultValue
}

func getEnvFloat64Value(key string, defaultValue float64) float64 {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatal(err)
		}

		return f
	}
	return defaultValue
}
