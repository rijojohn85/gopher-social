package main

import (
	"net/http"
)

// Health godoc
//
//	@Summary		Check health of api
//	@Description	check healtth of api
//	@Tags			health
//	@Produce		json
//	@Param			Authorization	header		string	true	"Auth payload"
//	@Success		201				{object}	map[string]string
//	@Failure		401				{object}	error
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [GET]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
