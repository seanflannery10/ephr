package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/helpers"
)

type connectionConfig struct {
	port int
	env  string
}

type smtpConfig struct {
	host     string
	port     int
	username string
	password string
	sender   string
}

type config struct {
	connection connectionConfig
	smtp       smtpConfig
	db         database.Config
}

func parseConfig() config {
	config := config{
		connection: connectionConfig{
			port: int(getEnvInt32Value("PORT", 4000)),
			env:  getEnvStrValue("ENV", "dev"),
		},
		smtp: smtpConfig{
			host:     getEnvStrValue("SMTP_HOST", ""),
			port:     int(getEnvInt32Value("SMTP_PORT", 25)),
			username: getEnvStrValue("SMTP_USERNAME", ""),
			password: getEnvStrValue("SMTP_PASSWORD", ""),
			sender:   getEnvStrValue("SMTP_SENDER", "Greenlight <no-reply@testdomain.com>"),
		},
		db: database.Config{
			DSN:                   getEnvStrValue("DB_URL", ""),
			MaxConns:              getEnvInt32Value("DB_MAX_CONNS", 25),
			MinConns:              getEnvInt32Value("DB_MIN_CONNS", 25),
			MaxConnLifetime:       getEnvStrValue("DB_MAX_CONN_LIFETIME", "15m"),
			MaxConnLifetimeJitter: getEnvStrValue("DB_MAX_CONN_LIFETIME_JITTER", "15m"),
			MaxConnIdleTime:       getEnvStrValue("DB_MAX_CONN_IDLE_TIME", "60m"),
		},
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	return config
}

func getEnvStrValue(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return defaultValue
}

// func getEnvBoolValue(key string, defaultValue bool) bool {
//	if value, ok := os.LookupEnv(key); ok && value != "" {
//		switch value {
//		case "true", "True", "1":
//			return true
//		default:
//			return false
//		}
//	}
//
//	return defaultValue
//}

func getEnvInt32Value(key string, defaultValue int32) int32 {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		return int32(i)
	}

	return defaultValue
}

// func getEnvFloat64Value(key string, defaultValue float64) float64 {
//	if value, ok := os.LookupEnv(key); ok && value != "" {
//		f, err := strconv.ParseFloat(value, 64)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		return f
//	}
//
//	return defaultValue
//}
