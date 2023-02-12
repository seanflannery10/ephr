package main

import (
	"context"
	"net/http"
	"time"

	"github.com/seanflannery10/ephr/internal/data"
	"github.com/seanflannery10/ephr/internal/queries"
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

	createUserParams := queries.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Activated:    false,
	}

	if data.ValidateNewUserParams(v, createUserParams); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	ctxCreateUser, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	user, err := app.queries.CreateUser(ctxCreateUser, createUserParams)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			v.AddError("email", "a user with this email address already exists")
			httperrors.FailedValidation(w, r, v)
		default:
			httperrors.ServerError(w, r, err)
			return
		}
	}

	addPermissionsForUserParms := queries.AddPermissionsForUserParams{
		UserID: user.ID,
		Code:   "movies:read",
	}

	ctxAddPermissionsForUser, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	_, err = app.queries.AddPermissionsForUser(ctxAddPermissionsForUser, addPermissionsForUserParms)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	createTokenParams, err := data.GenCreateTokenParams(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	ctxCreateToken, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	token, err := app.queries.CreateToken(ctxCreateToken, createTokenParams)
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
