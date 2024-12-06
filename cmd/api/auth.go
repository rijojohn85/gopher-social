package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/mailer"
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
	domain := env.GetString("DOMAIN", "http://localhost:8080")
	activationURL := fmt.Sprintf("%s/activate/%s", domain, plainToken)
	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	err = app.mailer.Send(
		mailer.UserWelcomeTemplate,
		user.Username,
		user.Email,
		vars,
		!isProdEnv,
	)
	if err != nil {
		app.logger.Errorw("Error sending mail", "error", err)
		//rollback user creation if email fails (SAGA Pattern)
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("Deleting email failed", "error", err)
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.logger.Errorw("Deleting email failed", "error", err)
				app.internalServerError(w, r, errors.New("user not found, contact developer"))
			default:
				app.internalServerError(w, r, err)
			}
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, plainToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
