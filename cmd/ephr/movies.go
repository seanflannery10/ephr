package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/jsonutil"
	"github.com/seanflannery10/ossa/read"
	"github.com/seanflannery10/ossa/validator"
	"net/http"
)

var (
	ctx               = context.Background()
	errRecordNotFound = errors.New("record not found")
	errEditConflict   = errors.New("edit conflict")
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := jsonutil.Read(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	movie := data.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := &validator.Validator{}

	if app.validateCreateMovie(v, movie); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	createdMovie, err := app.queries.CreateMovie(ctx, movie)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", createdMovie.ID))

	err = jsonutil.WriteWithHeaders(w, http.StatusCreated, map[string]any{"movie": createdMovie}, headers)
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := read.IDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	movie, err := app.queries.GetMovie(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errRecordNotFound):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}
		return
	}

	err = jsonutil.Write(w, http.StatusOK, map[string]any{"movie": movie})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := read.IDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	movie, err := app.queries.GetMovie(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errRecordNotFound):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}
		return
	}

	var input struct {
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = jsonutil.Read(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := &validator.Validator{}

	if app.validateMovie(v, movie); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	_, err = app.queries.UpdateMovie(ctx, movie)
	if err != nil {
		switch {
		case errors.Is(err, errEditConflict):
			httperrors.EditConflict(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}
		return
	}

	err = jsonutil.Write(w, http.StatusOK, map[string]any{"movie": movie})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
