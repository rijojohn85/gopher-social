package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//get the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, errors.New("missing authorization header"))
				return
			}

			//parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedError(w, r, errors.New("auth header format must be Basic {token}"))
				return
			}
			//decode it
			b, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}
			//check credentials
			user := app.config.auth.basic.user
			pass := app.config.auth.basic.pass
			creds := strings.SplitN(string(b), ":", 2)
			if len(creds) != 2 || creds[0] != user || creds[1] != pass {
				app.unauthorizedError(w, r, errors.New("invalid credentials"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
