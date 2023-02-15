package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		Genres:  pgtype.Array[string]{Elements: input.Genres},
	}

	v := validator.New()

	if data.ValidateCreateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.CreateMovie(r.Context(), params)
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

	movie, err := app.queries.GetMovie(r.Context(), id)
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
		params.Genres = pgtype.Array[string]{Elements: input.Genres}
	}

	v := validator.New()

	if data.ValidateUpdateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movie, err := app.queries.UpdateMovie(r.Context(), params)
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

	err = app.queries.DeleteMovie(r.Context(), id)
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
		Genres pgtype.Array[string]
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = helpers.ReadStringParam(qs, "title", "")
	input.Genres = pgtype.Array[string]{Elements: helpers.ReadCSVParam(qs, "genres", []string{})}

	input.Filters.Page = helpers.ReadIntParam(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadIntParam(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadStringParam(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	movies, metadata, err := app.queries.GetAllMoviesWithMetadata(r.Context(), input.Title, input.Genres, input.Filters)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"movies": movies, "metadata": metadata})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
