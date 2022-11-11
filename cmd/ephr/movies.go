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
		Title   *string   `json:"title"`
		Year    *int32    `json:"year"`
		Runtime *int32    `json:"runtime"`
		Genres  *[]string `json:"genres"`
	}

	err := jsonutil.Read(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	movieParams := data.CreateMovieParams{
		Title:   *input.Title,
		Year:    *input.Year,
		Runtime: *input.Runtime,
		Genres:  *input.Genres,
	}

	v := &validator.Validator{}

	if app.validateCreateMovie(v, movieParams); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.CreateMovie(ctx, movieParams)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = jsonutil.WriteWithHeaders(w, http.StatusCreated, map[string]any{"movie": movie}, headers)
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
	var input struct {
		Title   *string   `json:"title,omitempty"`
		Year    *int32    `json:"year,omitempty"`
		Runtime *int32    `json:"runtime,omitempty"`
		Genres  *[]string `json:"genres,omitempty"`
	}

	err := jsonutil.Read(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	id, err := read.IDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	movieParams := data.UpdateMovieParams{ID: id}

	if input.Title != nil {
		movieParams.UpdateTitle = true
		movieParams.Title = *input.Title
	}

	if input.Year != nil {
		movieParams.UpdateYear = true
		movieParams.Year = *input.Year
	}

	if input.Runtime != nil {
		movieParams.UpdateRuntime = true
		movieParams.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movieParams.UpdateGenres = true
		movieParams.Genres = *input.Genres
	}

	v := &validator.Validator{}

	//TODO Fix partial fields
	if app.validateUpdateMovie(v, movieParams); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.UpdateMovie(ctx, movieParams)
	if err != nil {
		switch {
		case errors.Is(err, errRecordNotFound):
			httperrors.NotFound(w, r)
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
