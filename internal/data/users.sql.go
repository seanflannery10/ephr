// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: users.sql

package data

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, name, email, password_hash, activated, version
`

type CreateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
	)
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
	Hash   []byte             `json:"hash"`
	Scope  string             `json:"scope"`
	Expiry pgtype.Timestamptz `json:"expiry"`
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
RETURNING id, created_at, name, email, password_hash, activated, version
`

type UpdateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
	ID           int64  `json:"id"`
	Version      int32  `json:"version"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
		arg.ID,
		arg.Version,
	)
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
