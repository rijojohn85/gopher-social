package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rijojohn85/social/internal/store"
)

type userKey string

const flrCtxKey userKey = "follower"

const userCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Activates/Registers a user
//	@Description	Activates/Registers a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token path		string true	"Invitation Token"
//	@Success		204		{string}	string "User Activated"
//	@Failure		404		{object}	error
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [PUT]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		app.badRequestError(w, r, errors.New("token required"))
		return
	}
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrInvalidToken):
			app.badRequestError(w, r, err)
		case errors.Is(err, store.ErrInvitationExpired):
			app.badRequestError(w, r, err)
		case errors.Is(err, store.ErrorNotFound):
			app.internalServerError(w, r, errors.New(
				"user not found. Contact developer"),
			)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, "user activated"); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [GET]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := userFromContext(r.Context())
	err := app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			param := chi.URLParam(r, "userID")

			userId, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				app.badRequestError(w, r, err)
				return
			}

			user := &store.User{}

			ctx := r.Context()
			err = app.store.Users.GetUser(ctx, user, userId)
			if err != nil {
				switch {
				case errors.Is(err, store.ErrorNotFound):
					app.notFoundError(w, r, err)
					return
				default:
					app.internalServerError(w, r, err)
					return
				}
			}
			ctx = context.WithValue(ctx, flrCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func userFromContext(ctx context.Context) *store.User {
	user := ctx.Value(flrCtxKey).(*store.User)
	return user
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		204		{string}	 string "User Followed"
//	@Failure		400		{object}	error "Bad Request: Payload missing/error"
//	@Failure		404		{object}	error "User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [PUT]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userCtxKey).(*store.User)
	followerUser := userFromContext(r.Context())
	userID := user.ID
	err := app.store.Users.AddFollower(r.Context(), userID, followerUser.ID)
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyFollows) {
			app.conflictRequestError(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user
//	@Description	Unfollows a user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		204		{string}	 string "User Unfollowed"
//	@Failure		400		{object}	error "Bad Request: Payload missing/error"
//	@Failure		404		{object}	error "User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [PUT]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := userFromContext(r.Context())
	user := r.Context().Value(userCtxKey).(*store.User)
	userID := user.ID
	err := app.store.Users.DeleteFollower(r.Context(), userID, followerUser.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
