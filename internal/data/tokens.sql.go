// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: tokens.sql

package data

import (
	"context"
	"time"
)

const createToken = `-- name: CreateToken :one
INSERT INTO tokens (hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4)
RETURNING hash
`

type CreateTokenParams struct {
	Hash   []byte
	UserID int64
	Expiry time.Time
	Scope  string
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) ([]byte, error) {
	row := q.db.QueryRow(ctx, createToken,
		arg.Hash,
		arg.UserID,
		arg.Expiry,
		arg.Scope,
	)
	var hash []byte
	err := row.Scan(&hash)
	return hash, err
}

const deleteAllTokensForUser = `-- name: DeleteAllTokensForUser :exec
DELETE
FROM tokens
WHERE scope = $1
  AND user_id = $2
`

type DeleteAllTokensForUserParams struct {
	Scope  string
	UserID int64
}

func (q *Queries) DeleteAllTokensForUser(ctx context.Context, arg DeleteAllTokensForUserParams) error {
	_, err := q.db.Exec(ctx, deleteAllTokensForUser, arg.Scope, arg.UserID)
	return err
}
