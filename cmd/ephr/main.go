package main

import (
	"expvar"
	"log"

	"github.com/seanflannery10/ephr/internal/queries"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/server"
)

type application struct {
	config  config
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

	app := &application{
		config:  cfg,
		queries: q,
	}

	srv := server.New(app.config.connection.port, app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
