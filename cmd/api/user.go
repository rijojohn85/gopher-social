package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rijojohn85/social/internal/store"
	"net/http"
	"strconv"
)

type userKey string

const userCtxKey userKey = "user"

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

			var user = &store.User{}

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

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := userFromContext(r.Context())
	//TODO: update after auth
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

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := userFromContext(r.Context())
	//TODO: update after auth
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
