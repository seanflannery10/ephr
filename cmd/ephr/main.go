package main

import (
	"fmt"
	"log"

	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/server"
)

type application struct {
	config  config
	queries *data.Queries
}

func main() {
	cfg := parseConfig()

	dbpool, err := database.New(cfg.db)
	if err != nil {
		log.Fatal(err, nil)
	}

	publishVars(dbpool)

	queries := data.New(dbpool)

	app := &application{
		config:  cfg,
		queries: queries,
	}

	address := fmt.Sprintf(":%d", app.config.connection.port)
	srv := server.New(address, app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
