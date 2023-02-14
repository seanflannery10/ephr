package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/ossa/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

func (q *Queries) NewToken(userID int64, ttl time.Duration, scope string) (Token, error) {
	params := CreateTokenParams{
		UserID: userID,
		Expiry: pgtype.Timestamptz{Time: time.Now().Add(ttl)},
		Scope:  scope,
	}

	token, err := q.CreateToken(context.Background(), params)
	if err != nil {
		return Token{}, err
	}

	randomBytes := make([]byte, 16)

	_, err = rand.Read(randomBytes)
	if err != nil {
		return Token{}, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
