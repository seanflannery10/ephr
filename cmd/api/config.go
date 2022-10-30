package main

import (
	"flag"
	"fmt"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/middleware"
	"os"
	"strings"
)

type config struct {
	port         int
	env          string
	db           database.Config
	authenticate middleware.AuthenticateConfig
	cors         middleware.CorsConfig
	rateLimit    middleware.RateLimitConfig
}

func parseFlags() config {
	var (
		cfg      config
		maxConns int
		minConns int
	)

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stag|prod)")

	flag.StringVar(&cfg.db.DSN, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&maxConns, "db-max-conns", 25, "PostgreSQL max connections")
	flag.IntVar(&minConns, "db-min-conns", 25, "PostgreSQL min connections")
	flag.StringVar(&cfg.db.MaxConnLifetime, "db-max-conn lifetime", "15m", "PostgreSQL max connection lifetime")
	flag.StringVar(&cfg.db.MaxConnLifetimeJitter, "db-max-conn-lifetime-jitter", "15m", "PostgreSQL max connection lifetime jitter")
	flag.StringVar(&cfg.db.MaxConnIdleTime, "db-max-conn-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.authenticate.JWKSURL, "auth-jwks-url", "", "URL of JWKS endpoint")
	flag.StringVar(&cfg.authenticate.APIURL, "auth-api-url", "", "URL of this API")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	flag.BoolVar(&cfg.rateLimit.Enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.rateLimit.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.rateLimit.Burst, "limiter-burst", 4, "Rate limiter maximum burst")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", ver)
		os.Exit(0)
	}

	cfg.db.MaxConns = int32(maxConns)
	cfg.db.MinConns = int32(minConns)

	return cfg
}
