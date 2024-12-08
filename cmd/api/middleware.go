package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rijojohn85/social/internal/store"
)

func (app *application) CheckPostOwernship(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userCtxKey).(*store.User)
		post := r.Context().Value(postCtx).(*store.Post)
		// if method is get or user is owner of post or method is patch and role is mod or
		// role is admin continue
		if (r.Method == http.MethodGet || post.UserID == user.ID) ||
			(r.Method == http.MethodPatch && user.RoleID == 2) ||
			(user.RoleID == 3) {
			next.ServeHTTP(w, r)
			return
		}
		app.forbiddenError(w, r, errors.New("unauthorizedError"))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.logger.Error("Basic Auth")
			// get the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, errors.New("missing authorization header"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedError(w, r, errors.New("auth header format must be Basic {token}"))
				return
			}
			// decode it
			b, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedError(w, r, err)
				return
			}
			// check credentials
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

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the auth header read and split
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, errors.New("missing authorization header"))
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, errors.New("auth header format must be Bearer {token}"))
			return
		}
		token := parts[1]

		// convert token to jwt token using validate token
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// get userID
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()
		user := &store.User{}
		err = app.store.Users.GetUser(ctx, user, int64(userID))
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtxKey, user)
		app.logger.Info("user logged", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
