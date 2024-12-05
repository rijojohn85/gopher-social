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

const userCtxKey userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fecthes a user profile by ID
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
			ctx = context.WithValue(ctx, userCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func userFromContext(ctx context.Context) *store.User {
	user := ctx.Value(userCtxKey).(*store.User)
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
	followerUser := userFromContext(r.Context())
	// TODO: update after auth
	userID := 1
	err := app.store.Users.AddFollower(r.Context(), int64(userID), followerUser.ID)
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
	// TODO: update after auth
	userID := 1
	err := app.store.Users.DeleteFollower(r.Context(), int64(userID), followerUser.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
