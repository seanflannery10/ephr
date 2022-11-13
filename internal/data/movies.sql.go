// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: movies.sql

package data

import (
	"context"
	"time"
)

const createMovie = `-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version
`

type CreateMovieParams struct {
	Title   string
	Year    int32
	Runtime int32
	Genres  []string
}

type CreateMovieRow struct {
	ID        int64
	CreatedAt time.Time
	Version   int32
}

func (q *Queries) CreateMovie(ctx context.Context, arg *CreateMovieParams) (*CreateMovieRow, error) {
	row := q.db.QueryRow(ctx, createMovie,
		arg.Title,
		arg.Year,
		arg.Runtime,
		arg.Genres,
	)
	var i CreateMovieRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Version)
	return &i, err
}

const deleteMovie = `-- name: DeleteMovie :exec
DELETE
FROM movies
WHERE id = $1
`

func (q *Queries) DeleteMovie(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteMovie, id)
	return err
}

const getAllMovies = `-- name: GetAllMovies :many
SELECT count(*) OVER (),
       id,
       created_at,
       title,
       year,
       runtime,
       genres,
       version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
  AND (genres @> $2 OR $2 = '{}')
ORDER BY CASE WHEN $3::bool THEN id END,
         CASE WHEN $4::bool THEN id END DESC,
         CASE WHEN $5::bool THEN title END,
         CASE WHEN $6::bool THEN title END DESC,
         CASE WHEN $7::bool THEN year END,
         CASE WHEN $8::bool THEN year END DESC,
         CASE WHEN $9::bool THEN runtime END,
         CASE WHEN $10::bool THEN runtime END DESC
OFFSET $11 LIMIT $12
`

type GetAllMoviesParams struct {
	Title       string
	Genres      []string
	IDAsc       bool
	IDDesc      bool
	TitleAsc    bool
	TitleDesc   bool
	YearAsc     bool
	YearDesc    bool
	RuntimeAsc  bool
	RuntimeDesc bool
	Offset      int32
	Limit       int32
}

type GetAllMoviesRow struct {
	Count     int64
	ID        int64
	CreatedAt time.Time
	Title     string
	Year      int32
	Runtime   int32
	Genres    []string
	Version   int32
}

func (q *Queries) GetAllMovies(ctx context.Context, arg *GetAllMoviesParams) ([]*GetAllMoviesRow, error) {
	rows, err := q.db.Query(ctx, getAllMovies,
		arg.Title,
		arg.Genres,
		arg.IDAsc,
		arg.IDDesc,
		arg.TitleAsc,
		arg.TitleDesc,
		arg.YearAsc,
		arg.YearDesc,
		arg.RuntimeAsc,
		arg.RuntimeDesc,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetAllMoviesRow
	for rows.Next() {
		var i GetAllMoviesRow
		if err := rows.Scan(
			&i.Count,
			&i.ID,
			&i.CreatedAt,
			&i.Title,
			&i.Year,
			&i.Runtime,
			&i.Genres,
			&i.Version,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMovie = `-- name: GetMovie :one
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1
`

func (q *Queries) GetMovie(ctx context.Context, id int64) (*Movie, error) {
	row := q.db.QueryRow(ctx, getMovie, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Title,
		&i.Year,
		&i.Runtime,
		&i.Genres,
		&i.Version,
	)
	return &i, err
}

const updateMovie = `-- name: UpdateMovie :one
UPDATE movies
SET title   = CASE WHEN $1::boolean THEN $2::TEXT ELSE title END,
    year    = CASE WHEN $3::boolean THEN $4::INTEGER ELSE year END,
    runtime = CASE WHEN $5::boolean THEN $6::INTEGER ELSE runtime END,
    genres  = CASE WHEN $7::boolean THEN $8::TEXT[] ELSE genres END,
    version = version + 1
WHERE id = $9
RETURNING version
`

type UpdateMovieParams struct {
	UpdateTitle   bool
	Title         string
	UpdateYear    bool
	Year          int32
	UpdateRuntime bool
	Runtime       int32
	UpdateGenres  bool
	Genres        []string
	ID            int64
}

func (q *Queries) UpdateMovie(ctx context.Context, arg *UpdateMovieParams) (int32, error) {
	row := q.db.QueryRow(ctx, updateMovie,
		arg.UpdateTitle,
		arg.Title,
		arg.UpdateYear,
		arg.Year,
		arg.UpdateRuntime,
		arg.Runtime,
		arg.UpdateGenres,
		arg.Genres,
		arg.ID,
	)
	var version int32
	err := row.Scan(&version)
	return version, err
}
