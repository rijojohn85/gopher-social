package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := NewTestApplication(t)
	mux := app.mount()
	testToken, _ := app.authenticator.GenerateToken(nil)
	t.Run("should not allow unauthenticated request", func(t *testing.T) {
		// check for 401 code
		req, err := http.NewRequest(
			http.MethodGet,
			"/v1/users/1",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
		w := executeRequest(t, req, mux)
		checkStatus(t, w.Code, http.StatusUnauthorized)
	})
	t.Run("should allow authenticated request", func(t *testing.T) {
		// check for 401 code
		req, err := http.NewRequest(
			http.MethodGet,
			"/v1/users/1",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(
			"Authorization",
			"Bearer "+testToken,
		)
		w := executeRequest(t, req, mux)
		checkStatus(t, w.Code, http.StatusOK)
	})
}
