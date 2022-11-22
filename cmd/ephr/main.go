package main

import (
	"flag"
	"fmt"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/database"
	"github.com/seanflannery10/ossa/server"
	"github.com/seanflannery10/ossa/version"
	"log"
	"os"
	"sync"
)

type application struct {
	config  config
	queries *data.Queries
	wg      sync.WaitGroup
}

func main() {
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version.Get())
		os.Exit(0)
	}

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

	address := fmt.Sprintf(":%d", app.config.port)
	srv := server.New(address, app.routes())

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
