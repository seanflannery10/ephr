-- migrate:up
CREATE TABLE IF NOT EXISTS movies
(
    id         bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title      text                        NOT NULL,
    year       integer                     NOT NULL,
    runtime    integer                     NOT NULL,
    genres     text[]                      NOT NULL,
    version    integer                     NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);
CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));

ALTER TABLE movies
    ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
ALTER TABLE movies
    ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies
    ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));

-- migrate:down
ALTER TABLE movies
    DROP CONSTRAINT IF EXISTS genres_length_check;
ALTER TABLE movies
    DROP CONSTRAINT IF EXISTS movies_runtime_check;
ALTER TABLE movies
    DROP CONSTRAINT IF EXISTS movies_year_check;

DROP INDEX IF EXISTS movies_genres_idx;
DROP INDEX IF EXISTS movies_title_idx;

DROP TABLE IF EXISTS movies;