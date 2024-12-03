package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) internalServerError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf("%v internal server error: %q path: %q error: %q", time.Now(), r.Method, r.URL.Path, err.Error())

	writeJsonError(
		w,
		http.StatusInternalServerError,
		"The server encountered an error",
	)
}

func (app *application) badRequestError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf(
		"%v Bad Request sent error: %q path: %q, body: %s error: %q",
		time.Now(),
		r.Method,
		r.URL.Path,
		r.Body,
		err.Error(),
	)

	writeJsonError(
		w,
		http.StatusBadRequest,
		err.Error(),
	)
}

func (app *application) conflictRequestError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf(
		"%v Conflict sent error: %q path: %q, body: %s error: %q",
		time.Now(),
		r.Method,
		r.URL.Path,
		r.Body,
		err.Error(),
	)

	writeJsonError(
		w, http.StatusConflict, err.Error(),
	)
}

func (app *application) notFoundError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	log.Printf(
		"%v Not Found error: %q path: %q, error: %q",
		time.Now(),
		r.Method,
		r.URL.Path,
		err.Error(),
	)

	writeJsonError(
		w,
		http.StatusNotFound,
		err.Error(),
	)
}
