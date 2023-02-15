package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/validator"
	"golang.org/x/exp/slog"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidatePasswordPlaintext(v, input.Password); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	hash, err := data.GetPasswordHash(input.Password)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	params := data.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Activated:    false,
	}

	if data.ValidateNewUserParams(v, params); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	user, err := app.queries.CreateUser(r.Context(), params)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			v.AddError("email", "a user with this email address already exists")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	_, err = app.queries.AddPermissionsForUser(r.Context(), data.AddPermissionsForUserParams{
		UserID: user.ID,
		Code:   "movies:read",
	})
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	token, err := app.queries.NewToken(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	app.server.Background(func() {
		input := map[string]any{
			"activationToken": token.Plaintext,
		}

		err = app.mailer.Send(user.Email, "token_activation.tmpl", input)
		if err != nil {
			slog.Error("email error", err)
		}
	})

	err = helpers.WriteJSON(w, http.StatusCreated, map[string]any{"authentication_token": token})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	user, err := app.queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  data.ScopeActivation,
		Expiry: pgtype.Timestamptz{Time: time.Now()},
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v.AddError("token", "invalid or expired activation token")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	user.Activated = true

	_, err = app.queries.UpdateUser(r.Context(), data.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Activated:    user.Activated,
		ID:           user.ID,
		Version:      user.Version,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			v.AddError("email", "a user with this email address already exists")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = app.queries.DeleteAllTokensForUser(r.Context(), data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"user": user})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	data.ValidatePasswordPlaintext(v, input.Password)
	data.ValidateTokenPlaintext(v, input.TokenPlaintext)

	if v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	user, err := app.queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  data.ScopePasswordReset,
		Expiry: pgtype.Timestamptz{Time: time.Now()},
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v.AddError("token", "invalid or expired password token")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	hash, err := data.GetPasswordHash(input.Password)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	_, err = app.queries.UpdateUser(r.Context(), data.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hash,
		Activated:    user.Activated,
		ID:           user.ID,
		Version:      user.Version,
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			v.AddError("email", "a user with this email address already exists")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = app.queries.DeleteAllTokensForUser(r.Context(), data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, map[string]any{"message": "your password was successfully reset"})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
