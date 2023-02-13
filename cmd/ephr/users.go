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

	paramsCreateUser := data.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Activated:    false,
	}

	if data.ValidateNewUserParams(v, paramsCreateUser); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	// Write user to db
	user, err := app.queries.CreateUser(r.Context(), paramsCreateUser)
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

	paramsAddPermissionsForUser := data.AddPermissionsForUserParams{
		UserID: user.ID,
		Code:   "movies:read",
	}

	// Write user permissions to db
	_, err = app.queries.AddPermissionsForUser(r.Context(), paramsAddPermissionsForUser)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	parmsCreateToken, _, err := data.GenCreateTokenParams(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	// Write token to db
	token, err := app.queries.CreateToken(r.Context(), parmsCreateToken)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	// TODO Fix
	// app.background(func() {
	//	data := map[string]any{
	//		"activationToken": token.Plaintext,
	//		"userID":          user.ID,
	//	}
	//
	//	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	//	if err != nil {
	//		app.logger.PrintError(err, nil)
	//	}
	// })

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

	paramsGetUserFromToken := data.GenGetUserFromTokenParams(input.TokenPlaintext, data.ScopeActivation)

	user, err := app.queries.GetUserFromToken(r.Context(), paramsGetUserFromToken)
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

	paramsUpdateUser := data.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Activated:    user.Activated,
		ID:           user.ID,
		Version:      user.Version,
	}

	_, err = app.queries.UpdateUser(r.Context(), paramsUpdateUser)
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

	paramsDeleteAllTokensForUser := data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	}

	err = app.queries.DeleteAllTokensForUser(r.Context(), paramsDeleteAllTokensForUser)
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

	paramsGetUserFromToken := data.GenGetUserFromTokenParams(input.TokenPlaintext, data.ScopePasswordReset)

	user, err := app.queries.GetUserFromToken(r.Context(), paramsGetUserFromToken)
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

	paramsUpdateUser := data.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hash,
		Activated:    user.Activated,
		ID:           user.ID,
		Version:      user.Version,
	}

	_, err = app.queries.UpdateUser(r.Context(), paramsUpdateUser)
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

	paramsDeleteAllTokensForUser := data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	}

	err = app.queries.DeleteAllTokensForUser(r.Context(), paramsDeleteAllTokensForUser)
	if err != nil {
		httperrors.ServerError(w, r, err)
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, map[string]any{"message": "your password was successfully reset"})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
