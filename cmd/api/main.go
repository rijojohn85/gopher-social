package main

import (
	"github.com/rijojohn85/social/internal/db"
	"github.com/rijojohn85/social/internal/db/auth"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/mailer"
	"github.com/rijojohn85/social/internal/store"
	"go.uber.org/zap"
	"time"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gopher

//	@contact.name	Rijo John
//	@contact.url	http://github.com/rijojohn85
//	@contact.email	rijojohn85@gmail.com

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString(
				"DB_ADDR",
				"postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable",
			),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME ", "15m"),
		},
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		mail: mailConfig{
			exp: time.Hour * 24,
			mailer: mailTripConfig{
				url:      env.GetString("MAILER_URL", "sandbox.smtp.mailtrap.io"),
				port:     env.GetInt("MAILER_PORT", 587),
				username: env.GetString("MAILER_USERNAME", ""),
				password: env.GetString("MAILER_PASSWORD", ""),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: "rijo",
				pass: "password",
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_SECRET", "hello_world"),
				aud:    env.GetString("TOKEN_AUD", "gophersocial"),
				exp:    time.Hour * 24 * 7,
			},
		},
	}
	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()
	database, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Panic(err.Error())
	}
	mailTripDialer := mailer.NewMailTripDialer(
		cfg.mail.mailer.url,
		cfg.mail.mailer.username,
		cfg.mail.mailer.password,
		env.GetString("MAILER_FROM_EMAIL", "rijojohn85@gmail.com"),
		cfg.mail.mailer.port,
	)
	logger.Info("Connected to Database pool")
	defer database.Close()
	storage := store.NewStorage(database)
	authetincator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.aud, cfg.auth.token.aud)
	app := &application{
		config:        cfg,
		store:         storage,
		logger:        logger,
		mailer:        mailTripDialer,
		authenticator: authetincator,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
