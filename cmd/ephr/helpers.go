package main

import (
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/validator"
	"math"
	"time"
)

func (app *application) validateCreateMovie(v *validator.Validator, createMovieParams data.CreateMovieParams) {
	v.CheckField(createMovieParams.Title != "", "title", "must be provided")
	v.CheckField(len(createMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.CheckField(createMovieParams.Year != 0, "year", "must be provided")
	v.CheckField(createMovieParams.Year >= 1888, "year", "must be greater than 1888")
	v.CheckField(createMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.CheckField(createMovieParams.Runtime != 0, "runtime", "must be provided")
	v.CheckField(createMovieParams.Runtime > 0, "runtime", "must be a positive integer")

	v.CheckField(createMovieParams.Genres != nil, "genres", "must be provided")
	v.CheckField(len(createMovieParams.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.CheckField(len(createMovieParams.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.CheckField(validator.NoDuplicates(createMovieParams.Genres), "genres", "must not contain duplicate values")
}

func (app *application) validateUpdateMovie(v *validator.Validator, updateMovieParams data.UpdateMovieParams) {
	if updateMovieParams.Title != "" {
		v.CheckField(len(updateMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")
	}

	if updateMovieParams.Year != 0 {
		v.CheckField(updateMovieParams.Year >= 1888, "year", "must be greater than 1888")
		v.CheckField(updateMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	}

	if updateMovieParams.Runtime != 0 {
		v.CheckField(updateMovieParams.Runtime > 0, "runtime", "must be a positive integer")
	}

	if updateMovieParams.Genres != nil {
		v.CheckField(len(updateMovieParams.Genres) >= 1, "genres", "must contain at least 1 genre")
		v.CheckField(len(updateMovieParams.Genres) <= 5, "genres", "must not contain more than 5 genres")
		v.CheckField(validator.NoDuplicates(updateMovieParams.Genres), "genres", "must not contain duplicate values")
	}
}

type Filters struct {
	Page     int
	PageSize int
}

func (f Filters) limit() int32 {
	return int32(f.PageSize)
}

func (f Filters) offset() int32 {
	return int32((f.Page - 1) * f.PageSize)
}

type Metadata struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	PageSize     int   `json:"page_size,omitempty"`
	FirstPage    int   `json:"first_page,omitempty"`
	LastPage     int   `json:"last_page,omitempty"`
	TotalRecords int64 `json:"total_records,omitempty"`
}

func (app *application) validateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page size must be greater than zero")
	v.Check(f.PageSize <= 100, "page size must be a maximum of 100")
}

func (app *application) calculateMetadata(totalRecords int64, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
