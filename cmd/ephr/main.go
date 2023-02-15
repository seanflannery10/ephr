package main

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ephr/internal/database"
	"github.com/seanflannery10/ephr/internal/mailer"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/server"
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

type application struct {
	config  config
	mailer  mailer.Mailer
	queries *data.Queries
	server  *server.Server
}

func main() {
	cfg := config{}

	err := envconfig.Process("ephr", cfg)
	if err != nil {
		log.Fatal(err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	m, err := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		log.Fatal(err, nil)
	}

	dbpool, err := database.New(cfg.db.dsn)
	if err != nil {
		log.Fatal(err, nil)
	}

	helpers.PublishCommonMetrics()
	expvar.Publish("database", expvar.Func(func() any {
		return dbpool.Stat()
	}))

	q := data.New(dbpool)

	app := &application{
		config:  cfg,
		queries: q,
		mailer:  m,
	}

	app.server = server.New(app.config.connection.port, app.routes())

	err = app.server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
