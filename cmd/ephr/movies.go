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

	params := data.CreateMovieParams{
		Title:   *input.Title,
		Year:    *input.Year,
		Runtime: *input.Runtime,
		Genres:  *input.Genres,
	}

	v := &validator.Validator{}

	if app.validateCreateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.CreateMovie(ctx, params)
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

	params := data.UpdateMovieParams{ID: id}

	if input.Title != nil {
		params.UpdateTitle = true
		params.Title = *input.Title
	}

	if input.Year != nil {
		params.UpdateYear = true
		params.Year = *input.Year
	}

	if input.Runtime != nil {
		params.UpdateRuntime = true
		params.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		params.UpdateGenres = true
		params.Genres = *input.Genres
	}

	v := &validator.Validator{}

	if app.validateUpdateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.UpdateMovie(ctx, params)
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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := read.IDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	err = app.queries.DeleteMovie(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errRecordNotFound):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}
		return
	}

	err = jsonutil.Write(w, http.StatusOK, map[string]any{"message": "movie successfully deleted"})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	PageSize     int   `json:"page_size,omitempty"`
	FirstPage    int   `json:"first_page,omitempty"`
	LastPage     int   `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string
		Genres []string
		Filters
	}

	v := &validator.Validator{}

	input.Title = read.String(r, "title", "")
	input.Genres = read.CSV(r, "genres", []string{})

	input.Filters.Page = read.Int(r, "page", 1, v)
	input.Filters.PageSize = read.Int(r, "page_size", 20, v)

	input.Filters.Sort = read.String(r, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	//if app.validateFilters(v, input.Filters); v.HasErrors() {
	//	httperrors.FailedValidation(w, r, v)
	//	return
	//}

	params := data.GetAllMoviesParams{}

	movies, err := app.queries.GetAllMovies(ctx, params)

	metadata := app.calculateMetadata(movies[0].Count, input.Filters.Page, input.Filters.PageSize)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	err = jsonutil.Write(w, http.StatusOK, map[string]any{"movies": movies, "metadata": metadata})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
