package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/mailer"
	"github.com/rijojohn85/social/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthPayload struct {
	Username string `json:"username" valid:"required,max=20"`
	Password string `json:"password" valid:"required,max=72,min=3"`
	Email    string `json:"email" valid:"required,email,max=255"`
	RoleID   int64  `json:"role_id" valid:"required, gte=1"`
}
type CreateTokenUserPayload struct {
	Email    string `json:"email" valid:"required,email,max=255"`
	Password string `json:"password" valid:"required,max=72,min=3"`
}

// RegisterUser godoc
//
//	@Summary		Register a user
//	@Description	Register a user and send email invite
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		AuthPayload	true	"Auth payload"
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
		RoleID:   payload.RoleID,
	}
	err = user.Password.Set(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// store the user
	plainToken := uuid.New().String()
	// encrpyt token
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
	activationURL := fmt.Sprintf("%s/v1/users/activate/%s", domain, plainToken)
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
		// rollback user creation if email fails (SAGA Pattern)
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("Deleting user failed", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusCreated, plainToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// Create Token godoc
//
//	@Summary		CreateToken
//	@Description	Create token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateTokenUserPayload	true	"CreateUserToken payload"
//	@Success		200		{object}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/token [POST]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload and validate
	var payload CreateTokenUserPayload

	err := readJson(w, r, &payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestError(w, r, err)
	}
	// fetch the user (check if the user exists) from payload
	user, err := app.store.Users.GetUserByEmail(r.Context(), payload.Email)
	if err != nil {
		app.unauthorizedError(w, r, err)
		return
	}
	// check credentials
	err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(payload.Password))
	if err != nil {
		app.unauthorizedError(w, r, errors.New("invalid credentials"))
		return
	}

	// generate token -> add claims
	claim := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.aud,
		"aud": app.config.auth.token.aud,
	}
	token, err := app.authenticator.GenerateToken(claim)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// send to client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
