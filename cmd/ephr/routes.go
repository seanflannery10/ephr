package main

import (
	"expvar"
	"github.com/seanflannery10/ossa/handlers"
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/httprouter"
	"github.com/seanflannery10/ossa/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	m := middleware.New()

	m.SetAuthenticateConfig(app.config.auth.JWKSURL, app.config.auth.APIURL)
	m.SetCorsConfig(app.config.cors.TrustedOrigins)
	m.SetRateLimitConfig(app.config.limit.Enabled, app.config.limit.RPS, app.config.limit.Burst)

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(httperrors.NotFound)
	router.MethodNotAllowed = http.HandlerFunc(httperrors.MethodNotAllowed)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handlers.Healthcheck)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	return m.Chain(m.Metrics, m.RecoverPanic, m.CORS, m.RateLimit, m.Authenticate).Then(router)
}
