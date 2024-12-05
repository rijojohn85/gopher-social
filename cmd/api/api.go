package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rijojohn85/social/docs" // this is required to generate swagger docs
	"github.com/rijojohn85/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	store  store.Storage
	config config
}

type config struct {
	addr   string
	env    string
	apiURL string
	db     dbConfig
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

	r.Route(
		"/v1", func(r chi.Router) {
			r.Get("/health", app.healthCheckHandler)
			docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
			r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

			r.Route(
				"/posts", func(r chi.Router) {
					r.Post("/", app.createPostHandler)

					r.Route(
						"/{postID}", func(r chi.Router) {
							r.Use(app.postContextMiddleware)
							r.Get("/", app.getPostHandler)
							r.Post("/comments", app.createCommentHandler)
							r.Patch("/", app.patchPostHandler)
							r.Delete("/", app.deletePostHandler)
						},
					)
				},
			)
			r.Route(
				"/users", func(r chi.Router) {
					r.Route(
						"/{userID}", func(r chi.Router) {
							r.Use(app.userContextMiddleware)
							r.Get("/", app.getUserHandler)
							r.Put("/follow", app.followUserHandler)
							r.Put("/unfollow", app.unfollowUserHandler)
						},
					)
					r.Group(func(r chi.Router) {
						r.Get("/feed", app.getUserFeedHandler)
					})
				},
			)
		},
	)

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
	log.Printf("server has started on http://localhost%s", srv.Addr)
	return srv.ListenAndServe()
}
