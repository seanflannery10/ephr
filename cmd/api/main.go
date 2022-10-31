package main

import (
	"fmt"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/log"
	"github.com/seanflannery10/ossa/server"
	"sync"
)

type application struct {
	config config
	wg     sync.WaitGroup
}

func main() {

	cfg := parseConfig()

	db, err := database.New(cfg.db)
	if err != nil {
		log.Fatal(err, nil)
	}

	publishVars(db)

	app := &application{
		config: cfg,
	}

	srv := server.New(fmt.Sprintf(":%d", app.config.port), app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err, nil)
	}
}
