package main

import (
	"log"

	"github.com/rijojohn85/social/internal/db"
	"github.com/rijojohn85/social/internal/env"
	"github.com/rijojohn85/social/internal/store"
)

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
		env: env.GetString("ENV", "development"),
	}
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Print("Connected to Database pool")
	defer db.Close()
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
