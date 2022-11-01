package main

import (
	"expvar"
	"github.com/seanflannery10/ossa/errors"
	"github.com/seanflannery10/ossa/handlers"
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

	router.NotFound = http.HandlerFunc(errors.NotFound)
	router.MethodNotAllowed = http.HandlerFunc(errors.MethodNotAllowed)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", handlers.Healthcheck)

	return m.Chain(m.Metrics, m.RecoverPanic, m.CORS, m.RateLimit, m.Authenticate).Then(router)
}
