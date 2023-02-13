package data

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/seanflannery10/ossa/validator"
	"golang.org/x/crypto/bcrypt"
)

func GenGetUserFromTokenParams(tokenPlaintext string, scope string) GetUserFromTokenParams {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	params := GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  scope,
		Expiry: time.Now(),
	}

	return params
}

func GetPasswordHash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ComparePasswords(plaintextPassword string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.RgxEmail), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateNewUserParams(v *validator.Validator, user CreateUserParams) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, fmt.Sprintf("%v", user.Email))
}
