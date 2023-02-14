package data

import (
	"time"

	"github.com/seanflannery10/ossa/validator"
)

func SetSort(sort string, params *GetAllMoviesParams) {
	switch sort {
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
}

func ValidateCreateMovie(v *validator.Validator, createMovieParams CreateMovieParams) {
	v.Check(createMovieParams.Title != "", "title", "must be provided")
	v.Check(len(createMovieParams.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(createMovieParams.Year != 0, "year", "must be provided")
	v.Check(createMovieParams.Year >= 1888, "year", "must be greater than 1888")
	v.Check(createMovieParams.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(createMovieParams.Runtime != 0, "runtime", "must be provided")
	v.Check(createMovieParams.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(createMovieParams.Genres.Elements != nil, "genres", "must be provided")
	v.Check(len(createMovieParams.Genres.Elements) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(createMovieParams.Genres.Elements) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.NoDuplicates(createMovieParams.Genres.Elements), "genres", "must not contain duplicate values")
}

func ValidateUpdateMovie(v *validator.Validator, updateMovieParams UpdateMovieParams) {
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

	if updateMovieParams.Genres.Elements != nil {
		v.Check(len(updateMovieParams.Genres.Elements) >= 1, "genres", "must contain at least 1 genre")
		v.Check(len(updateMovieParams.Genres.Elements) <= 5, "genres", "must not contain more than 5 genres")
		v.Check(validator.NoDuplicates(updateMovieParams.Genres.Elements), "genres", "must not contain duplicate values")
	}
}
