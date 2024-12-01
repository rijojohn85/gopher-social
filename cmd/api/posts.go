package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rijojohn85/social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		// TODO: change after auth
		UserID: 1,
		Tags:   payload.Tags,
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
	}
	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequestError(w, r, err)
	}

	ctx := r.Context()
	post := &store.Post{}
	err = app.store.Posts.GetPostById(ctx, post, int64(id))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
		}
	}
	comments, err := app.store.Comments.GetByPostID(ctx, int64(id))
	if err != nil {
		app.internalServerError(w, r, err)
	}
	post.Commennts = comments

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}
