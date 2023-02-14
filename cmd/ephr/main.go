package main

import (
	"expvar"
	"log"

	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ephr/internal/database"
	"github.com/seanflannery10/ephr/internal/mailer"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/server"
)

type application struct {
	config  config
	mailer  mailer.Mailer
	queries *data.Queries
	server  *server.Server
}

func main() {
	cfg := parseConfig()

	m, err := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		log.Fatal(err, nil)
	}

	dbpool, err := database.New(cfg.db)
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

	srv := server.New(app.config.connection.port, app.routes())

	app.server = srv

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
