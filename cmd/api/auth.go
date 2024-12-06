package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"github.com/rijojohn85/social/internal/store"
	"net/http"
)

type AuthPayload struct {
	Username string `json:"username" valid:"required,max=20"`
	Password string `json:"password" valid:"required,max=72,min=3"`
	Email    string `json:"email" valid:"required,email,max=255"`
}

// RegisterUser godoc
//
//	@Summary		Register a user
//	@Description	Register a user and send email invite
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		AuthPayload true	"Auth payload"
//	@Success		201		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/users [POST]
func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var payload AuthPayload
	err := readJson(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// store the user
	plainToken := uuid.New().String()
	//encrpyt token
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	err = app.store.Users.CreateAndInvite(
		r.Context(),
		user,
		hashToken,
		app.config.mail.exp,
	)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateUsername):
			app.conflictRequestError(w, r, err)
		case errors.Is(err, store.ErrDuplicateEmail):
			app.conflictRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	//TODO: send email
	if err := app.jsonResponse(w, http.StatusCreated, plainToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
