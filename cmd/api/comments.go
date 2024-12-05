package main

import (
	"fmt"
	"net/http"

	"github.com/rijojohn85/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=255"`
}

// CreateComment godoc
//
//	@Summary		Creates comment
//	@Description	Creates a comment with payload for a particular post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int						true	"postID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		200		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{pathID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
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
	post := getPostFromContext(r)
	if post == nil {
		app.badRequestError(w, r, fmt.Errorf("no post found in context"))
		return
	}
	comment := &store.Comment{
		// to be updated after auth
		UserId:  1,
		Content: payload.Content,
		PostID:  post.ID,
	}
	err = app.store.Comments.Create(r.Context(), comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
