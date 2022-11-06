package main

import (
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/validator"
	"time"
)

func (app *application) validateMovie(v *validator.Validator, movie data.CreateMovieParams) {
	v.CheckField(movie.Title != "", "title", "must be provided")
	v.CheckField(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.CheckField(movie.Year != 0, "year", "must be provided")
	v.CheckField(movie.Year >= 1888, "year", "must be greater than 1888")
	v.CheckField(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.CheckField(movie.Runtime != 0, "runtime", "must be provided")
	v.CheckField(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.CheckField(movie.Genres != nil, "genres", "must be provided")
	v.CheckField(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.CheckField(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.CheckField(validator.NoDuplicates(movie.Genres), "genres", "must not contain duplicate values")
}
