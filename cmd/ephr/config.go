package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
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

	err := envconfig.Process("ephr", config)
	if err != nil {
		log.Fatal(err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	return config
}
