package main

import (
	"expvar"
	"github.com/seanflannery10/ephr/internal/mailer"
	"log"

	"github.com/seanflannery10/ephr/internal/queries"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/server"
)

type application struct {
	config  config
	mailer  mailer.Mailer
	queries *queries.Queries
}

func main() {
	cfg := parseConfig()

	dbpool, err := database.New(cfg.db)
	if err != nil {
		log.Fatal(err, nil)
	}

	helpers.PublishCommonMetrics()
	expvar.Publish("database", expvar.Func(func() any {
		return dbpool.Stat()
	}))

	q := queries.New(dbpool)

	m, err := mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		log.Fatal(err, nil)
	}

	app := &application{
		config:  cfg,
		queries: q,
		mailer:  m,
	}

	srv := server.New(app.config.connection.port, app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
