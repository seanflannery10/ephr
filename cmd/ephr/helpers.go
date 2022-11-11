package main

import (
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/validator"
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
	v.CheckField(updateMovieParams.Title != "", "title", "must be provided")
	v.CheckField(len(updateMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.CheckField(updateMovieParams.Year != 0, "year", "must be provided")
	v.CheckField(updateMovieParams.Year >= 1888, "year", "must be greater than 1888")
	v.CheckField(updateMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.CheckField(updateMovieParams.Runtime != 0, "runtime", "must be provided")
	v.CheckField(updateMovieParams.Runtime > 0, "runtime", "must be a positive integer")

	v.CheckField(updateMovieParams.Genres != nil, "genres", "must be provided")
	v.CheckField(len(updateMovieParams.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.CheckField(len(updateMovieParams.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.CheckField(validator.NoDuplicates(updateMovieParams.Genres), "genres", "must not contain duplicate values")
}
