-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version;

-- name: DeleteMovie :exec
DELETE
FROM movies
WHERE id = $1;

-- name: GetMovie :one
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1;

-- name: GetAllMovies :many
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
LIMIT $3 OFFSET $4;

-- name: UpdateMovie :one
UPDATE movies
SET title   = $1,
    year    = $2,
    runtime = $3,
    genres  = $4,
    version = version + 1
WHERE id = $5
  AND version = $6
RETURNING version;