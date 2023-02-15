package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/seanflannery10/ossa/helpers"
)

type config struct {
	connection struct {
		port int
		env  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	db struct {
		dsn string
	}
}

func parseConfig() config {
	config := config{}

	config.connection.port = int(getEnvInt32Value("PORT", 4000))
	config.connection.env = getEnvStrValue("ENV", "dev")

	config.smtp.host = getEnvStrValue("SMTP_HOST", "smtp.mailtrap.io")
	config.smtp.port = int(getEnvInt32Value("SMTP_PORT", 25))
	config.smtp.username = getEnvStrValue("SMTP_USERNAME", "")
	config.smtp.password = getEnvStrValue("SMTP_PASSWORD", "")
	config.smtp.sender = getEnvStrValue("SMTP_SENDER", "Greenlight <no-reply@testdomain.com>")

	config.db.dsn = getEnvStrValue("DB_DSN", "")

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
