package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rijojohn85/social/internal/db/auth"
	"github.com/rijojohn85/social/internal/store"
	"github.com/rijojohn85/social/internal/store/cache"
	"go.uber.org/zap"
)

func NewTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockCache()
	mockAuth := auth.NewTestAuthenticator("hello world")

	app := application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: mockAuth,
	}
	return &app
}

func executeRequest(t *testing.T, req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	return w
}

func checkStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got StatusCode %d, wanted StatusCode %d", got, want)
	}
}
