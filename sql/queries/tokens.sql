-- name: CreateToken :one
INSERT INTO tokens (plaintext, hash, user_id, expiry, scope)
VALUES ('', $1, $2, $3, $4)
RETURNING *;

-- name: DeleteAllTokensForUser :exec
DELETE
FROM tokens
WHERE scope = $1
  AND user_id = $2;