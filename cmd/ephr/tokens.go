package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ossa/helpers"
	"github.com/seanflannery10/ossa/httperrors"
	"github.com/seanflannery10/ossa/validator"
	"golang.org/x/exp/slog"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	user, err := app.queries.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			httperrors.InvalidAuthenticationToken(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	match, err := data.ComparePasswords(input.Password, user.PasswordHash)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	if !match {
		httperrors.InvalidAuthenticationToken(w, r)
		return
	}

	token, err := app.queries.NewToken(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, map[string]any{"authentication_token": token})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	user, err := app.queries.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			httperrors.InvalidAuthenticationToken(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if !user.Activated {
		v.AddError("email", "user account must be activated")
		httperrors.FailedValidation(w, r, v)

		return
	}

	token, err := app.queries.NewToken(r.Context(), user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	app.server.Background(func() {
		input := map[string]any{
			"passwordResetToken": token.Plaintext,
		}

		err = app.mailer.Send(user.Email, "token_password_reset.tmpl", input)
		if err != nil {
			slog.Error("email error", err)
		}
	})

	msg := "an email will be sent to you containing password reset instructions"

	err = helpers.WriteJSON(w, http.StatusAccepted, map[string]any{"message": msg})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	user, err := app.queries.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			httperrors.InvalidAuthenticationToken(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if user.Activated {
		v.AddError("email", "user has already been activated")
		httperrors.FailedValidation(w, r, v)

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

	msg := "an email will be sent to you containing activation instructions"

	err = helpers.WriteJSON(w, http.StatusAccepted, map[string]any{"message": msg})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
