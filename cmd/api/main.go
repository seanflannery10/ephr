package main

import (
	"fmt"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/log"
	"github.com/seanflannery10/ossa/server"
	"sync"
)

type application struct {
	config  config
	queries *data.Queries
	wg      sync.WaitGroup
}

func main() {
	cfg := parseConfig()

	db, err := database.New(cfg.db)
	if err != nil {
		log.Fatal(err, nil)
	}

	publishVars(db)

	queries := data.New(db)

	app := &application{
		config:  cfg,
		queries: queries,
	}

	address := fmt.Sprintf(":%d", app.config.port)
	srv := server.New(address, app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err, nil)
	}
}
