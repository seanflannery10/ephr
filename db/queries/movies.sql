-- name: CreateMovie :one
INSERT INTO movies (name, bio)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteMovie :exec
DELETE
FROM movies
WHERE id = $1;

-- name: GetMovie :one
SELECT *
FROM movies
WHERE id = $1
LIMIT 1;

-- name: ListMovies :many
SELECT *
FROM movies
ORDER BY name;

-- name: UpdateMovie :one
UPDATE movies
set name = $2,
    bio  = $3
WHERE id = $1
RETURNING *;