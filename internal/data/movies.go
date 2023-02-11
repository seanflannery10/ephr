package data

import (
	"time"

	"github.com/seanflannery10/ephr/internal/queries"
	"github.com/seanflannery10/ossa/validator"
)

func ValidateCreateMovie(v *validator.Validator, createMovieParams queries.CreateMovieParams) {
	v.Check(createMovieParams.Title != "", "title", "must be provided")
	v.Check(len(createMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(createMovieParams.Year != 0, "year", "must be provided")
	v.Check(createMovieParams.Year >= 1888, "year", "must be greater than 1888")
	v.Check(createMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(createMovieParams.Runtime != 0, "runtime", "must be provided")
	v.Check(createMovieParams.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(createMovieParams.Genres != nil, "genres", "must be provided")
	v.Check(len(createMovieParams.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(createMovieParams.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.NoDuplicates(createMovieParams.Genres), "genres", "must not contain duplicate values")
}

func ValidateUpdateMovie(v *validator.Validator, updateMovieParams queries.UpdateMovieParams) {
	if updateMovieParams.Title != "" {
		v.Check(len(updateMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")
	}

	if updateMovieParams.Year != 0 {
		v.Check(updateMovieParams.Year >= 1888, "year", "must be greater than 1888")
		v.Check(updateMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	}

	if updateMovieParams.Runtime != 0 {
		v.Check(updateMovieParams.Runtime > 0, "runtime", "must be a positive integer")
	}

	if updateMovieParams.Genres != nil {
		v.Check(len(updateMovieParams.Genres) >= 1, "genres", "must contain at least 1 genre")
		v.Check(len(updateMovieParams.Genres) <= 5, "genres", "must not contain more than 5 genres")
		v.Check(validator.NoDuplicates(updateMovieParams.Genres), "genres", "must not contain duplicate values")
	}
}
