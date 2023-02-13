// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: users.sql

package queries

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version
`

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash []byte
	Activated    bool
}

type CreateUserRow struct {
	ID        int64
	CreatedAt time.Time
	Version   int32
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
	)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Version)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const getUserFromToken = `-- name: GetUserFromToken :one
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
         INNER JOIN tokens
                    ON users.id = tokens.user_id
WHERE tokens.hash = $1
  AND tokens.scope = $2
  AND tokens.expiry > $3
`

type GetUserFromTokenParams struct {
	Hash   []byte
	Scope  string
	Expiry time.Time
}

func (q *Queries) GetUserFromToken(ctx context.Context, arg GetUserFromTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserFromToken, arg.Hash, arg.Scope, arg.Expiry)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name          = $1,
    email         = $2,
    password_hash = $3,
    activated     = $4,
    version       = version + 1
WHERE id = $5
  AND version = $6
RETURNING version
`

type UpdateUserParams struct {
	Name         string
	Email        string
	PasswordHash []byte
	Activated    bool
	ID           int64
	Version      int32
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (int32, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
		arg.ID,
		arg.Version,
	)
	var version int32
	err := row.Scan(&version)
	return version, err
}
