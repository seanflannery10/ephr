-- migrate:up
CREATE TABLE movies
(
    id   BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    bio  text
);

-- migrate:down
DROP TABLE movies;
