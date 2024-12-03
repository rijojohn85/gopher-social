package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rijojohn85/social/internal/store"
	"net/http"
	"strconv"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

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

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
