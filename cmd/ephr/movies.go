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
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/jsonutil"
	"github.com/seanflannery10/ossa/read"
	"github.com/seanflannery10/ossa/validator"
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

	params := &data.CreateMovieParams{
		Title:   *input.Title,
		Year:    *input.Year,
		Runtime: *input.Runtime,
		Genres:  *input.Genres,
	}

	v := validator.New()

	if app.validateCreateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.CreateMovie(ctx, params) //nolint:contextcheck
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.GetMovie(ctx, id) //nolint:contextcheck
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
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

	params := &data.UpdateMovieParams{ID: id}

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

	v := validator.New()

	if app.validateUpdateMovie(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	movie, err := app.queries.UpdateMovie(ctx, params) //nolint:contextcheck
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = app.queries.DeleteMovie(ctx, id) //nolint:contextcheck
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
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

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	title := read.String(r, "title", "")
	genres := read.CSV(r, "genres", []string{})

	v := validator.New()
	filters := Filters{}

	filters.Page = read.Int(r, "page", 1, v)
	filters.PageSize = read.Int(r, "page_size", 20, v)

	filters.Sort = read.String(r, "sort", "id")
	filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if app.validateFilters(v, filters); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	params := &data.GetAllMoviesParams{
		Title:  title,
		Genres: genres,
		Offset: filters.offset(),
		Limit:  filters.limit(),
	}

	switch filters.Sort {
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	movies, err := app.queries.GetAllMovies(ctx, params) //nolint:contextcheck
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	count, err := app.queries.GetMovieCount(ctx) //nolint:contextcheck
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	metadata := app.calculateMetadata(count, filters.Page, filters.PageSize)

	err = jsonutil.Write(w, http.StatusOK, map[string]any{"movies": movies, "metadata": metadata})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
