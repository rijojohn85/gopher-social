package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rijojohn85/social/internal/db/auth"
	"github.com/rijojohn85/social/internal/mailer"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rijojohn85/social/docs" // this is required to generate swagger docs
	"github.com/rijojohn85/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	store         store.Storage
	mailer        mailer.Client
	authenticator auth.Authenticator
	logger        *zap.SugaredLogger
	config        config
}

type config struct {
	mail   mailConfig
	addr   string
	env    string
	apiURL string
	auth   authConfig
	db     dbConfig
}
type authConfig struct {
	basic basicAuthConfig
	token tokenConfig
}
type tokenConfig struct {
	secret string
	aud    string
	exp    time.Duration
}
type basicAuthConfig struct {
	user string
	pass string
}
type mailConfig struct {
	mailer mailTripConfig
	exp    time.Duration
}
type mailTripConfig struct {
	url      string
	username string
	password string
	port     int
}
type dbConfig struct {
	addr         string
	maxIdleTime  string
	maxOpenConns int
	maxIdleConns int
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	/*  Set a timeout value on the request context (ctx), that will signal
		* through ctx.Done() that the request has timed out and further
		* processing should be stopped.
	* */
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Post("/", app.createPostHandler)
			r.Route(
				"/{postID}", func(r chi.Router) {
					r.Use(app.postContextMiddleware)
					r.Get("/", app.getPostHandler)
					r.Post("/comments", app.createCommentHandler)
					r.With(app.CheckPostOwernship).Patch("/", app.patchPostHandler)
					r.With(app.CheckPostOwernship).Delete("/", app.deletePostHandler)
				})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/users", app.registerUser)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	srv := http.Server{
		Addr:    app.config.addr,
		Handler: mux, WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}
	app.logger.Infow(
		"server has started",
		"addr",
		app.config.addr,
		"env",
		app.config.env,
	)
	return srv.ListenAndServe()
}
