package main

import (
	"github.com/rijojohn85/social/internal/store"
	"net/http"
)

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
