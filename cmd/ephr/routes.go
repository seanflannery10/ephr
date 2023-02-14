package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/seanflannery10/ossa/handlers"
	"github.com/seanflannery10/ossa/httperrors"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	r.NotFound(httperrors.NotFound)
	r.MethodNotAllowed(httperrors.MethodNotAllowed)

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	r.Get("/v1/healthcheck", handlers.Healthcheck)

	r.Route("/v1/movies", func(r chi.Router) {
		r.Get("/", app.listMoviesHandler)
		r.Post("/", app.createMovieHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.showMovieHandler)
			r.Patch("/", app.updateMovieHandler)
			r.Delete("/", app.deleteMovieHandler)
		})
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", app.registerUserHandler)
		r.Put("/activated", app.activateUserHandler)
		r.Put("/password", app.updateUserPasswordHandler)
	})

	r.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/authentication", app.createAuthenticationTokenHandler)
		r.Put("/activation", app.createActivationTokenHandler)
		r.Put("/password-reset", app.createPasswordResetTokenHandler)
	})

	return r
}
