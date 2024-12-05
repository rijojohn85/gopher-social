package main

import (
	"github.com/rijojohn85/social/internal/db"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/store"
	"go.uber.org/zap"
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
	logger.Info("Connected to Database pool")
	defer database.Close()
	storage := store.NewStorage(database)
	app := &application{
		config: cfg,
		store:  storage,
		logger: logger,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
