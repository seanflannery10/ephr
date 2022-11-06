package main

import (
	"context"
	"fmt"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/errors"
	"github.com/seanflannery10/ossa/json"
	"github.com/seanflannery10/ossa/validator"
	"net/http"
)

var ctx = context.Background()

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := json.Decode(w, r, &input)
	if err != nil {
		errors.BadRequest(w, r, err)
		return
	}

	movie := data.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := &validator.Validator{}

	if app.validateMovie(v, movie); v.HasErrors() {
		errors.FailedValidation(w, r, v)
		return
	}

	createdMovie, err := app.queries.CreateMovie(ctx, movie)
	if err != nil {
		errors.ServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", createdMovie.ID))

	err = json.EncodeWithHeaders(w, http.StatusCreated, map[string]any{"movie": createdMovie}, headers)
	if err != nil {
		errors.ServerError(w, r, err)
	}
}

//func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
//	id, err := read.IDParam(r)
//	if err != nil {
//		errors.NotFound(w, r)
//		return
//	}
//
//	movie, err := app.models.Movies.Get(id)
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrRecordNotFound):
//			errors.NotFound(w, r)
//		default:
//			errors.ServerError(w, r, err)
//		}
//		return
//	}
//
//	err = json.Encode(w, http.StatusOK, map[string]any{"movie": movie})
//	if err != nil {
//		errors.ServerError(w, r, err)
//	}
//}
