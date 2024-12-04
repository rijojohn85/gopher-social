package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := app.store.Posts.GetUserFeed(ctx, int64(5))
	if err != nil {
		app.internalServerError(w, r, err)
	}
	err = app.jsonResponse(w, http.StatusOK, posts)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
