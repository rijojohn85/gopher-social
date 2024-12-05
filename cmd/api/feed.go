package main

import (
	"net/http"

	"github.com/rijojohn85/social/internal/store"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.UserFeed
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort
	ctx := r.Context()
	var fq store.PaginatedFeedQuery
	fq.Limit = 10
	fq.Offset = 0
	fq.Sort = "desc"
	fq, err := fq.Parse(r)
	posts, err := app.store.Posts.GetUserFeed(ctx, int64(5), fq)
	if err != nil {
		app.internalServerError(w, r, err)
	}
	err = app.jsonResponse(w, http.StatusOK, posts)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
