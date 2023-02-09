package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	params := data.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if app.validateCreateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.CreateMovie(ctx, params)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = helpers.WriteJSONWithHeaders(w, http.StatusCreated, map[string]any{"movie": movie}, headers)
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.GetMovie(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"movie": movie})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   *string  `json:"title,omitempty"`
		Year    *int32   `json:"year,omitempty"`
		Runtime *int32   `json:"runtime,omitempty"`
		Genres  []string `json:"genres,omitempty"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	id, err := helpers.ReadIDParam(r)
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
		params.Genres = input.Genres
	}

	v := validator.New()

	if app.validateUpdateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.UpdateMovie(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
			httperrors.EditConflict(w, r)
			return
		}
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"movie": movie})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err = app.queries.DeleteMovie(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"message": "movie successfully deleted"})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = helpers.ReadStringParam(qs, "title", "")
	input.Genres = helpers.ReadCSVParam(qs, "genres", []string{})

	input.Filters.Page = helpers.ReadIntParam(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadIntParam(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadStringParam(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if app.validateFilters(v, input.Filters); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	params := data.GetAllMoviesParams{
		Title:  input.Title,
		Genres: input.Genres,
		Offset: input.Filters.offset(),
		Limit:  input.Filters.limit(),
	}

	switch input.Filters.Sort {
	case "id":
		params.IDAsc = true
	case "-id":
		params.IDDesc = true
	case "title":
		params.TitleAsc = true
	case "-title":
		params.TitleDesc = true
	case "year":
		params.YearAsc = true
	case "-year":
		params.YearDesc = true
	case "runtime":
		params.RuntimeAsc = true
	case "-runtime":
		params.RuntimeDesc = true
	default:
		params.IDAsc = true
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	movies, err := app.queries.GetAllMovies(ctx, params)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	count, err := app.queries.GetMovieCount(ctx)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	metadata := app.calculateMetadata(count, input.Filters.Page, input.Filters.PageSize)

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"movies": movies, "metadata": metadata})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
