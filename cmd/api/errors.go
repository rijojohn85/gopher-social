package main

import (
	"net/http"
	"time"

	"github.com/rijojohn85/social/internal/store"
)

func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"Forbibben Access Error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
		"accessingUser",
		r.Context().Value(userCtxKey).(*store.User),
	)

	errJson := writeJsonError(
		w,
		http.StatusForbidden,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}

func (app *application) internalServerError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Errorw(
		"Internal Server Error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
	)

	errJson := writeJsonError(
		w,
		http.StatusInternalServerError,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}

func (app *application) badRequestError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Warnw(
		"Bad Request error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
	)

	errJson := writeJsonError(
		w,
		http.StatusBadRequest,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}

func (app *application) conflictRequestError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Errorw(
		"Conflict Request error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
	)

	errJson := writeJsonError(
		w,
		http.StatusConflict,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}

func (app *application) notFoundError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	app.logger.Warnw(
		"NotFound Request error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
	)
	errJson := writeJsonError(
		w,
		http.StatusNotFound,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw(
		"Unauthorized error",
		"path",
		r.URL.Path,
		"error",
		err.Error(),
		"method",
		r.Method,
		"time",
		time.Now(),
	)
	errJson := writeJsonError(
		w,
		http.StatusUnauthorized,
		err.Error(),
	)
	if errJson != nil {
		app.internalServerError(w, r, errJson)
		return
	}
}
