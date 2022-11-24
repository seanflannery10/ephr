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
SELECT id,
       created_at,
       title,
       year,
       runtime,
       genres,
       version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', @title) OR @title = '')
  AND (genres @> @genres OR @genres = '{}')
ORDER BY CASE WHEN @id_asc::bool THEN id END,
         CASE WHEN @id_desc::bool THEN id END DESC,
         CASE WHEN @title_asc::bool THEN title END,
         CASE WHEN @title_desc::bool THEN title END DESC,
         CASE WHEN @year_asc::bool THEN year END,
         CASE WHEN @year_desc::bool THEN year END DESC,
         CASE WHEN @runtime_asc::bool THEN runtime END,
         CASE WHEN @runtime_desc::bool THEN runtime END DESC
OFFSET @offset_ LIMIT @limit_;

-- name: GetMovieCount :one
SELECT count(*)
FROM movies;

-- name: UpdateMovie :one
UPDATE movies
SET title   = CASE WHEN @update_title::boolean THEN @title::TEXT ELSE title END,
    year    = CASE WHEN @update_year::boolean THEN @year::INTEGER ELSE year END,
    runtime = CASE WHEN @update_runtime::boolean THEN @runtime::INTEGER ELSE runtime END,
    genres  = CASE WHEN @update_genres::boolean THEN @genres::TEXT[] ELSE genres END,
    version = version + 1
WHERE id = @id
RETURNING title, year, runtime, genres, version;