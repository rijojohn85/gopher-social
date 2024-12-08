package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rijojohn85/social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   string   `json:"title" validate:"omitempty,max=100"`
	Content string   `json:"content" validate:"omitempty,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates post
//	@Description	Creates a post with payload
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJson(
		w,
		r,
		&payload,
	); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	user := r.Context().Value(userCtxKey).(*store.User)
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  user.ID,
		Tags:    payload.Tags,
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(
		ctx,
		post,
	); err != nil {
		app.internalServerError(w, r, err)
	}
	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UpdatePost godoc
//
//	@Summary		Updates post
//	@Description	Updates a post with payload
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"postID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [patch]
func (app *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload UpdatePostPayload
	post := getPostFromContext(r)

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
	}
	if payload.Content != "" {
		post.Content = payload.Content
	}
	if payload.Tags != nil {
		post.Tags = payload.Tags
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}
	ctx := r.Context()
	if err := app.store.Posts.Update(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes post
//	@Description	Deletes a post with id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"postID"
//	@Success		200		{string}	string	"post deleted"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [DELETE]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(
		r,
		"postID",
	)
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequestError(w, r, err)
	}
	ctx := r.Context()
	if err := app.store.Posts.Delete(
		ctx,
		int64(id),
	); err != nil {
		if errors.Is(
			err,
			store.ErrorNotFound,
		) {
			app.notFoundError(w, r, err)
			return
		} else {
			app.internalServerError(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// GetPost godoc
//
//	@Summary		Gets post
//	@Description	Gets a post with id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int			true	"postID"
//	@Success		200		{object}	store.Post	"post deleted"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [GET]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post := getPostFromContext(r)
	comments, err := app.store.Comments.GetByPostID(
		ctx,
		post.ID,
	)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			idParam := chi.URLParam(r, "postID")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				app.badRequestError(w, r, err)
				return
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
					return
				}
			}
			ctx = context.WithValue(r.Context(), postCtx, post)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func getPostFromContext(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
