package main

import (
	"github.com/seanflannery10/ossa/httprouter"
	"github.com/seanflannery10/ossa/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	m := middleware.New()

	return m.RecoverPanic(router)
}
