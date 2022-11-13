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
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', @title) OR @title = '')
  AND (genres @> @genres OR @genres = '{}')
-- TODO order by any filed and pick direction
ORDER BY id
OFFSET @offset_ LIMIT @limit_;

-- name: UpdateMovie :one
UPDATE movies
SET title   = CASE WHEN @update_title::boolean THEN @title::TEXT ELSE title END,
    year    = CASE WHEN @update_year::boolean THEN @year::INTEGER ELSE year END,
    runtime = CASE WHEN @update_runtime::boolean THEN @runtime::INTEGER ELSE runtime END,
    genres  = CASE WHEN @update_genres::boolean THEN @genres::TEXT[] ELSE genres END,
    version = version + 1
WHERE id = @id
RETURNING version;